package ec2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

var client *ec2.Client
var ctx = context.TODO()

func Configure(cfg aws.Config) {
	client = ec2.NewFromConfig(cfg)
}

func DeleteDefaultVpc() {
	if client == nil {
		log.Fatalf("ec2.Client is nil")
	}
	describeVpcsOutput := describeVpcs()
	defaultVpcId := describeVpcsOutput.Vpcs[0].VpcId

	deleteSubnets(defaultVpcId)

	deleteInternetGateways(defaultVpcId)

	_, deleteVpcErr := client.DeleteVpc(ctx, &ec2.DeleteVpcInput{
		VpcId: defaultVpcId,
	})
	if deleteVpcErr != nil {
		log.Fatalf("ec2.DeleteVpc failed because %v", deleteVpcErr)
	}
}

func createFilter(name *string, value *string) []types.Filter {
	return []types.Filter{
		{
			Name:   name,
			Values: []string{*value},
		},
	}
}

func deleteInternetGateway(vpcId *string, igwId *string) {
	detachInternetGatewayInput := ec2.DetachInternetGatewayInput{
		InternetGatewayId: igwId,
		VpcId:             vpcId,
	}
	_, detachInternetGatewayErr := client.DetachInternetGateway(
		ctx, &detachInternetGatewayInput,
	)
	if detachInternetGatewayErr != nil {
		log.Fatalf("ec2.DetachInternetGateway failed because %v", detachInternetGatewayErr)
	}

	deleteInternetGatewayInput := ec2.DeleteInternetGatewayInput{
		InternetGatewayId: igwId,
	}
	_, deleteInternetGatewayErr := client.DeleteInternetGateway(
		ctx, &deleteInternetGatewayInput,
	)
	if deleteInternetGatewayErr != nil {
		log.Fatalf("ec2.DeleteInternetGateway failed because %v", deleteInternetGatewayErr)
	}
}

func deleteInternetGateways(vpcId *string) {
	describeInternetGatewaysInput := ec2.DescribeInternetGatewaysInput{
		Filters: createFilter(aws.String("attachment.vpc-id"), vpcId),
	}
	describeInternetGatewaysOutput, describeInternetGatewaysErr := client.DescribeInternetGateways(
		ctx, &describeInternetGatewaysInput,
	)
	if describeInternetGatewaysErr != nil {
		log.Fatalf("ec2.DescribeInternetGateways failed because %v", describeInternetGatewaysErr)
	}

	for _, igw := range describeInternetGatewaysOutput.InternetGateways {
		deleteInternetGateway(vpcId, igw.InternetGatewayId)
	}
}

func deleteSubnets(vpcId *string) {
	describeSubnetsOutput, describeSubnetsErr := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: createFilter(aws.String("vpc-id"), vpcId),
	})
	if describeSubnetsErr != nil {
		log.Fatalf("ec2.DescribeSubnets failed because %v", describeSubnetsErr)
	}
	for _, subnet := range describeSubnetsOutput.Subnets {
		if vpcId == subnet.VpcId {
			_, deleteSubnetErr := client.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
				SubnetId: subnet.SubnetId,
			})
			if deleteSubnetErr != nil {
				log.Fatalf("ec2.DeleteSubnet %v failed because %v", &subnet.SubnetId, deleteSubnetErr)
			}
		}
	}
}

func describeVpcs() *ec2.DescribeVpcsOutput {
	describeVpcsOutput, describeVpcsErr := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		Filters: createFilter(aws.String("is-default"), aws.String("true")),
	})
	if describeVpcsErr != nil {
		log.Fatalf("ec2.DescribeVpcs failed because %v", describeVpcsErr)
	}
	return describeVpcsOutput
}
