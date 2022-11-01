package main

import (
	"account-vending-machine/config"
	"account-vending-machine/service/ec2"
	"account-vending-machine/service/organizations"
	"account-vending-machine/service/sts"
	"account-vending-machine/types"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event types.Event) (string, error) {
	cfg := config.Get()
	role := aws.String("OrganizationAccountAccessRole")

	organizations.Configure(cfg)
	createAccountOutput := organizations.CreateAccount(event, role)

	sts.Configure(cfg)
	assumeRoleCfg := sts.AssumeRole(createAccountOutput, role)

	ec2.Configure(assumeRoleCfg)
	ec2.DeleteDefaultVpc()

	return "Success", nil
}
