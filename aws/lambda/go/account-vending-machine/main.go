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
	ouRootId := ""
	ouName := ""
	role := "OrganizationAccountAccessRole"

	organizationsClient := organizations.NewFromConfig(cfg)
	createAccountOutput := createAccount(organizationsClient, account, domain, environment, role)
	accountId := aws.ToString(createAccountOutput.CreateAccountStatus.AccountId)

	updateAccountOrganizationalUnit(organizationsClient, accountId, ouRootId, ouName)

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
	describeVpcsOutput, describeVpcsErr := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{
		Filters: []ec2Types.Filter{
			{
				Name: aws.String("is-default"),
				Values: []string{
					"true",
				},
			},
		},
	})
	if describeVpcsErr != nil {
		log.Fatalf("ec2.DescribeVpcs failed because %v", describeVpcsErr)
	}
	defaultVpcId := describeVpcsOutput.Vpcs[0].VpcId

	deleteSubnets(client, defaultVpcId)

	deleteInternetGateways(client, defaultVpcId)

	_, deleteVpcErr := client.DeleteVpc(context.TODO(), &ec2.DeleteVpcInput{
		VpcId: defaultVpcId,
	})
	if deleteVpcErr != nil {
		log.Fatalf("ec2.DeleteVpc failed because %v", deleteVpcErr)
	}
}

func deleteInternetGateway(client *ec2.Client, vpcId *string, igwId *string) {
	detachInternetGatewayInput := ec2.DetachInternetGatewayInput{
		InternetGatewayId: igwId,
		VpcId:             vpcId,
	}
	_, detachInternetGatewayErr := client.DetachInternetGateway(
		context.TODO(), &detachInternetGatewayInput,
	)
	if detachInternetGatewayErr != nil {
		log.Fatalf("ec2.DetachInternetGateway failed because %v", detachInternetGatewayErr)
	}

	deleteInternetGatewayInput := ec2.DeleteInternetGatewayInput{
		InternetGatewayId: igwId,
	}
	_, deleteInternetGatewayErr := client.DeleteInternetGateway(
		context.TODO(), &deleteInternetGatewayInput,
	)
	if deleteInternetGatewayErr != nil {
		log.Fatalf("ec2.DeleteInternetGateway failed because %v", deleteInternetGatewayErr)
	}
}

func deleteInternetGateways(client *ec2.Client, vpcId *string) {
	describeInternetGatewaysInput := ec2.DescribeInternetGatewaysInput{
		Filters: []ec2Types.Filter{
			{
				Name: aws.String("attachment.vpc-id"),
				Values: []string{
					aws.ToString(vpcId),
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
		deleteInternetGateway(client, vpcId, igw.InternetGatewayId)
	}
}

func deleteSubnets(client *ec2.Client, vpcId *string) {
	describeSubnetsOutput, describeSubnetsErr := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []string{
					aws.ToString(vpcId),
				},
			},
		},
	})
	if describeSubnetsErr != nil {
		log.Fatalf("ec2.DescribeSubnets failed because %v", describeSubnetsErr)
	}
	for _, subnet := range describeSubnetsOutput.Subnets {
		if vpcId == subnet.VpcId {
			_, deleteSubnetErr := client.DeleteSubnet(context.TODO(), &ec2.DeleteSubnetInput{
				SubnetId: subnet.SubnetId,
			})
			if deleteSubnetErr != nil {
				log.Fatalf("ec2.DeleteSubnet %v failed because %v", &subnet.SubnetId, deleteSubnetErr)
			}
		}
	}
}

func updateAccountOrganizationalUnit(
	client *organizations.Client,
	accountId string,
	ouRootId string,
	ouName string,
) {
	listOrganizationalUnitsForParentInput := organizations.ListOrganizationalUnitsForParentInput{
		ParentId: &ouRootId,
	}
	listOrganizationalUnitsForParentOutput, listOrganizationalUnitsForParentErr := client.ListOrganizationalUnitsForParent(
		context.TODO(), &listOrganizationalUnitsForParentInput,
	)
	if listOrganizationalUnitsForParentErr != nil {
		log.Fatalf("organizations.ListOrganizationalUnitsForParent failed because %v", listOrganizationalUnitsForParentErr)
	}

	var ouId *string
	for _, ou := range listOrganizationalUnitsForParentOutput.OrganizationalUnits {
		if &ouName == ou.Name {
			ouId = ou.Id
			break
		}
	}
	if ouId == nil {
		log.Fatalf("Could not find 'ouId' for '%v' with 'organizationsRootId': %v", ouName, ouRootId)
	}

	_, moveAccountErr := client.MoveAccount(context.TODO(), &organizations.MoveAccountInput{
		AccountId:           &accountId,
		DestinationParentId: &ouRootId,
		SourceParentId:      ouId,
	})
	if moveAccountErr != nil {
		log.Fatalf("organizations.MoveAccount failed because %v", moveAccountErr)
	}
}
