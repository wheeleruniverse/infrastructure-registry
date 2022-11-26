package main

import (
	"account-vending-machine/config"
	"account-vending-machine/service/cloudformation"
	"account-vending-machine/service/ec2"
	"account-vending-machine/service/organizations"
	"account-vending-machine/service/sts"
	"account-vending-machine/types"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
)

func main() {
	lambda.Start(Handler)
}

func Handler(request types.Request) (types.Response, error) {
	if "Create" == request.RequestType {
		create(request)
	}
	return types.Response{
		LogicalResourceId:  request.LogicalResourceId,
		PhysicalResourceId: uuid.New().String(),
		RequestId:          request.RequestId,
		StackId:            request.StackId,
		Status:             "SUCCESS",
	}, nil
}

func create(request types.Request) {
	cfg := config.Get()
	role := aws.String("OrganizationAccountAccessRole")

	organizations.Configure(cfg)
	createAccountOutput := organizations.CreateAccount(request.ResourceProperties, role)

	sts.Configure(cfg)
	assumeRoleCfg := sts.AssumeRole(createAccountOutput, role)

	cloudformation.Configure(assumeRoleCfg)
	cloudformation.CreateBaseline(request.ResourceProperties)

	ec2.Configure(assumeRoleCfg)
	ec2.DeleteDefaultVpc()
}
