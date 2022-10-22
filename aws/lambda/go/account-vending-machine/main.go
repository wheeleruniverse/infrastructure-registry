package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"log"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() (string, error) {
	cfg, cfgErr := config.LoadDefaultConfig(context.TODO())
	if cfgErr != nil {
		log.Fatalf("LoadDefaultConfig failed because %v", cfgErr)
	}

	organizationsClient := organizations.NewFromConfig(cfg)
	createAccount(organizationsClient, "", "", "")

	// TODO move account to OU

	// TODO assume role for new account

	// TODO with assumed role delete default vpc

	return "Success", nil
}

func createAccount(
	organizationsClient *organizations.Client,
	account string,
	domain string,
	environment string,
) *organizations.CreateAccountOutput {
	email := aws.String(account + "@wheelerswebservices.com")
	createAccountInput := organizations.CreateAccountInput{
		AccountName:            &account,
		Email:                  email,
		IamUserAccessToBilling: types.IAMUserAccessToBillingDeny,
		RoleName:               aws.String("OrganizationAccountAccessRole"),
		Tags: []types.Tag{
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
	createAccountOutput, createAccountErr := organizationsClient.CreateAccount(context.TODO(), &createAccountInput)
	if createAccountErr != nil {
		log.Fatalf("organizations.CreateAccount failed because %v", createAccountErr)
	}
	return createAccountOutput
}