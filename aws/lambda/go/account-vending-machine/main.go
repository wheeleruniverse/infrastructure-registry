package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationsTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() (string, error) {
	cfgRegion := config.WithRegion("us-east-1")
	cfg, cfgErr := config.LoadDefaultConfig(context.TODO(), cfgRegion)
	if cfgErr != nil {
		log.Fatalf("cfg: LoadDefaultConfig failed because %v", cfgErr)
	}

	account := ""
	domain := ""
	environment := ""
	role := "OrganizationAccountAccessRole"

	organizationsClient := organizations.NewFromConfig(cfg)
	createAccountOutput := createAccount(organizationsClient, account, domain, environment, role)
	accountId := aws.ToString(createAccountOutput.CreateAccountStatus.AccountId)

	// TODO move account to OU

	stsClient := sts.NewFromConfig(cfg)
	assumeRoleOutput := assumeRole(stsClient, accountId, role)
	assumeRoleCredentials := assumeRoleOutput.Credentials

	assumeRoleCfg, assumeRoleCfgErr := config.LoadDefaultConfig(
		context.TODO(),
		cfgRegion,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     *assumeRoleCredentials.AccessKeyId,
				SecretAccessKey: *assumeRoleCredentials.SecretAccessKey,
				SessionToken:    *assumeRoleCredentials.SessionToken,
				Source:          "assumeRoleCredentials",
			},
		}),
	)
	if assumeRoleCfgErr != nil {
		log.Fatalf("assumeRoleCfg: LoadDefaultConfig failed because %v", assumeRoleCfgErr)
	}

	ec2Client := ec2.NewFromConfig(assumeRoleCfg)
	deleteDefaultVpc(ec2Client)

	return "Success", nil
}

func assumeRole(
	client *sts.Client,
	accountId string,
	role string,
) *sts.AssumeRoleOutput {
	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::" + accountId + ":role/" + role),
		RoleSessionName: aws.String("NewAccountRole"),
		DurationSeconds: aws.Int32(900),
	}
	assumeRoleOutput, assumeRoleErr := client.AssumeRole(context.TODO(), &assumeRoleInput)
	if assumeRoleErr != nil {
		log.Fatalf("sts.AssumeRole failed because %v", assumeRoleErr)
	}
	return assumeRoleOutput
}

func createAccount(
	client *organizations.Client,
	account string,
	domain string,
	environment string,
	role string,
) *organizations.CreateAccountOutput {
	email := aws.String(account + "@wheelerswebservices.com")
	createAccountInput := organizations.CreateAccountInput{
		AccountName:            &account,
		Email:                  email,
		IamUserAccessToBilling: organizationsTypes.IAMUserAccessToBillingDeny,
		RoleName:               &role,
		Tags: []organizationsTypes.Tag{
			{
				Key:   aws.String("Account"),
				Value: &account,
			},
			{
				Key:   aws.String("Environment"),
				Value: &environment,
			},
			{
				Key:   aws.String("Domain"),
				Value: &domain,
			},
			{
				Key:   aws.String("Owner"),
				Value: email,
			},
		},
	}
	createAccountOutput, createAccountErr := client.CreateAccount(context.TODO(), &createAccountInput)
	if createAccountErr != nil {
		log.Fatalf("organizations.CreateAccount failed because %v", createAccountErr)
	}
	return createAccountOutput
}

func deleteDefaultVpc(client *ec2.Client) {

	// find vpcs
	describeVpcsOutput, describeVpcsErr := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	if describeVpcsErr != nil {
		log.Fatalf("ec2.DescribeVpcs failed because %v", describeVpcsErr)
	}

	// find default vpc
	var defaultVpcId *string
	for _, vpc := range describeVpcsOutput.Vpcs {
		if ec2Types.TenancyDefault == vpc.InstanceTenancy {
			defaultVpcId = vpc.VpcId
			break
		}
	}
	if defaultVpcId == nil {
		// no default vpc
		return
	}

	// find subnets
	describeSubnetsOutput, describeSubnetsErr := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []string{
					aws.ToString(defaultVpcId),
				},
			},
		},
	})
	if describeSubnetsErr != nil {
		log.Fatalf("ec2.DescribeSubnets failed because %v", describeSubnetsErr)
	}

	// delete subnets
	for _, subnet := range describeSubnetsOutput.Subnets {
		if defaultVpcId == subnet.VpcId {
			_, deleteSubnetErr := client.DeleteSubnet(context.TODO(), &ec2.DeleteSubnetInput{
				SubnetId: subnet.SubnetId,
			})
			if deleteSubnetErr != nil {
				log.Fatalf("ec2.DeleteSubnet %v failed because %v", &subnet.SubnetId, deleteSubnetErr)
			}
		}
	}

	// describe internet gateways
	describeInternetGatewaysInput := ec2.DescribeInternetGatewaysInput{
		Filters: []ec2Types.Filter{
			{
				Name: aws.String("attachment.vpc-id"),
				Values: []string{
					aws.ToString(defaultVpcId),
				},
			},
		},
	}
	describeInternetGatewaysOutput, describeInternetGatewaysErr := client.DescribeInternetGateways(
		context.TODO(), &describeInternetGatewaysInput,
	)
	if describeInternetGatewaysErr != nil {
		log.Fatalf("ec2.DescribeInternetGateways failed because %v", describeInternetGatewaysErr)
	}

	for _, igw := range describeInternetGatewaysOutput.InternetGateways {
		// detach internet gateway
		detachInternetGatewayInput := ec2.DetachInternetGatewayInput{
			InternetGatewayId: igw.InternetGatewayId,
			VpcId:             defaultVpcId,
		}
		_, detachInternetGatewayErr := client.DetachInternetGateway(
			context.TODO(), &detachInternetGatewayInput,
		)
		if detachInternetGatewayErr != nil {
			log.Fatalf("ec2.DetachInternetGateway failed because %v", detachInternetGatewayErr)
		}

		// delete internet gateway
		deleteInternetGatewayInput := ec2.DeleteInternetGatewayInput{
			InternetGatewayId: igw.InternetGatewayId,
		}
		_, deleteInternetGatewayErr := client.DeleteInternetGateway(
			context.TODO(), &deleteInternetGatewayInput,
		)
		if deleteInternetGatewayErr != nil {
			log.Fatalf("ec2.DeleteInternetGateway failed because %v", deleteInternetGatewayErr)
		}
	}

	// delete vpc
	_, deleteVpcErr := client.DeleteVpc(context.TODO(), &ec2.DeleteVpcInput{
		VpcId: defaultVpcId,
	})
	if deleteVpcErr != nil {
		log.Fatalf("ec2.DeleteVpc failed because %v", deleteVpcErr)
	}
}

/*

   igw_response = ec2_client.describe_internet_gateways()
   for i in range(0,len(igw_response['InternetGateways'])):
       for j in range(0,len(igw_response['InternetGateways'][i]['Attachments'])):
           if(igw_response['InternetGateways'][i]['Attachments'][j]['VpcId'] == default_vpcid):
               default_igw = igw_response['InternetGateways'][i]['InternetGatewayId']
   #print(default_igw)
   detach_default_igw_response = ec2_client.detach_internet_gateway(InternetGatewayId=default_igw,VpcId=default_vpcid,DryRun=False)
   delete_internet_gateway_response = ec2_client.delete_internet_gateway(InternetGatewayId=default_igw)

   #print("Default IGW " + currentregion + "Deleted.")

   time.sleep(10)
   delete_vpc_response = ec2_client.delete_vpc(VpcId=default_vpcid,DryRun=False)
   print("Deleted Default VPC in {}".format(currentregion))
   return delete_vpc_response
*/
