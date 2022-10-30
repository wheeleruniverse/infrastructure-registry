package main

import (
	"account-vending-machine/service/ec2"
	"account-vending-machine/service/organizations"
	"account-vending-machine/service/sts"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

	role := aws.String("OrganizationAccountAccessRole")

	organizations.Configure(cfg)
	createAccountOutput := organizations.CreateAccount(organizations.CreateAccountInput{
		AccountName: aws.String(""),
		Domain:      aws.String(""),
		Environment: aws.String(""),
		Role:        role,
		OuRootId:    aws.String(""),
		OuName:      aws.String(""),
	})

	sts.Configure(cfg)
	assumeRoleCfg := sts.AssumeRole(sts.AssumeRoleInput{
		AccountId: createAccountOutput.AccountId,
		Role:      role,
	})

	ec2.Configure(assumeRoleCfg)
	ec2.DeleteDefaultVpc()

	return "Success", nil
}
