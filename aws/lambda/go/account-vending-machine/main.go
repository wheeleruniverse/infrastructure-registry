package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

	account := ""
	domain := ""
	environment := ""
	role := "OrganizationAccountAccessRole"

	organizationsClient := organizations.NewFromConfig(cfg)
	createAccountOutput := createAccount(organizationsClient, account, domain, environment, role)
	accountId := aws.ToString(createAccountOutput.CreateAccountStatus.AccountId)

	// TODO move account to OU

	stsClient := sts.NewFromConfig(cfg)
	assumeRole(stsClient, accountId, role)

	// TODO with assumed role delete default vpc

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
		IamUserAccessToBilling: types.IAMUserAccessToBillingDeny,
		RoleName:               &role,
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
	createAccountOutput, createAccountErr := client.CreateAccount(context.TODO(), &createAccountInput)
	if createAccountErr != nil {
		log.Fatalf("organizations.CreateAccount failed because %v", createAccountErr)
	}
	return createAccountOutput
}
