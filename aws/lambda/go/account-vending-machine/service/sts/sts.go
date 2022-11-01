package sts

import (
	"account-vending-machine/config"
	"account-vending-machine/types"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"
)

var client *sts.Client
var ctx = context.TODO()

func Configure(cfg aws.Config) {
	client = sts.NewFromConfig(cfg)
}

func AssumeRole(account types.Account, role *string) aws.Config {
	if client == nil {
		log.Fatalf("sts.Client is nil")
	}
	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::" + *account.AccountId + ":role/" + *role),
		RoleSessionName: aws.String("NewAccountRole"),
		DurationSeconds: aws.Int32(900),
	}
	assumeRoleOutput, assumeRoleErr := client.AssumeRole(ctx, &assumeRoleInput)
	if assumeRoleErr != nil {
		log.Fatalf("sts.AssumeRole failed because %v", assumeRoleErr)
	}
	assumeRoleCredentials := assumeRoleOutput.Credentials
	return config.GetWithStaticCredentials(assumeRoleCredentials)
}
