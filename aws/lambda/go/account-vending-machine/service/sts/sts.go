package sts

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"log"
)

var client *sts.Client
var ctx = context.TODO()
var region string

func Configure(cfg aws.Config) {
	client = sts.NewFromConfig(cfg)
	region = cfg.Region
}

type AssumeRoleInput struct {
	AccountId *string
	Role      *string
}

func AssumeRole(input AssumeRoleInput) aws.Config {
	if client == nil {
		log.Fatalf("sts.Client is nil")
	}
	accountId := *input.AccountId
	role := *input.Role
	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::" + accountId + ":role/" + role),
		RoleSessionName: aws.String("NewAccountRole"),
		DurationSeconds: aws.Int32(900),
	}
	assumeRoleOutput, assumeRoleErr := client.AssumeRole(ctx, &assumeRoleInput)
	if assumeRoleErr != nil {
		log.Fatalf("sts.AssumeRole failed because %v", assumeRoleErr)
	}
	assumeRoleCredentials := assumeRoleOutput.Credentials
	return retrieveConfigWithAssumedRoleCredentials(assumeRoleCredentials)
}

func retrieveConfigWithAssumedRoleCredentials(assumeRoleCredentials *types.Credentials) aws.Config {
	assumeRoleCfg, assumeRoleCfgErr := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     *assumeRoleCredentials.AccessKeyId,
				SecretAccessKey: *assumeRoleCredentials.SecretAccessKey,
				SessionToken:    *assumeRoleCredentials.SessionToken,
				Source:          "assumeRoleCredentials",
			},
		}),
		config.WithRegion(region),
	)
	if assumeRoleCfgErr != nil {
		log.Fatalf("LoadDefaultConfig failed because %v", assumeRoleCfgErr)
	}
	return assumeRoleCfg
}
