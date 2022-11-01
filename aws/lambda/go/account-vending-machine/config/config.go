package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"log"
)

var ctx = context.TODO()
var region = config.WithRegion("us-east-1")

func Get() aws.Config {
	return loadDefaultConfig(region)
}

func GetWithStaticCredentials(creds *types.Credentials) aws.Config {
	provider := config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     *creds.AccessKeyId,
			SecretAccessKey: *creds.SecretAccessKey,
			SessionToken:    *creds.SessionToken,
			Source:          "StaticCredentials",
		},
	})
	return loadDefaultConfig(provider, region)
}

func loadDefaultConfig(optFns ...func(*config.LoadOptions) error) aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		log.Fatalf("config.LoadDefaultConfig failed because %v", err)
	}
	return cfg
}
