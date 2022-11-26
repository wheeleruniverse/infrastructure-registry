package cloudformation

import (
	"account-vending-machine/types"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"log"
)

var client *cloudformation.Client
var ctx = context.TODO()

func Configure(cfg aws.Config) {
	client = cloudformation.NewFromConfig(cfg)
}

func CreateBaseline(event types.Event) {
	if client == nil {
		log.Fatalf("cloudformation.Client is nil")
	}
	createStackInput := cloudformation.CreateStackInput{
		Capabilities:     createCapabilities(),
		OnFailure:        cloudformationTypes.OnFailureRollback,
		Parameters:       createParameters(event),
		RoleARN:          aws.String("arn:aws:iam::028079862371:role/wheelerswebservices-infrastructure-prd-baseline"),
		StackName:        aws.String(event.AccountName + "-baseline-stack"),
		Tags:             createTags(event),
		TemplateURL:      aws.String("https://wheelerswebservices-infrastructure-prd-cf-templates.s3.amazonaws.com/baseline.yml"),
		TimeoutInMinutes: aws.Int32(60),
	}
	_, createStackErr := client.CreateStack(ctx, &createStackInput)
	if createStackErr != nil {
		log.Fatalf("cloudformation.CreateStack failed because %v", createStackErr)
	}
}

func createCapabilities() []cloudformationTypes.Capability {
	return []cloudformationTypes.Capability{
		cloudformationTypes.CapabilityCapabilityNamedIam,
	}
}

func createParameters(event types.Event) []cloudformationTypes.Parameter {
	return []cloudformationTypes.Parameter{
		{
			ParameterKey:   aws.String("pAccount"),
			ParameterValue: &event.AccountName,
		},
		{
			ParameterKey:   aws.String("pDomain"),
			ParameterValue: &event.Domain,
		},
		{
			ParameterKey:   aws.String("pEnvironment"),
			ParameterValue: &event.Environment,
		},
		{
			ParameterKey:   aws.String("pOwner"),
			ParameterValue: &event.Owner,
		},
		{
			ParameterKey:   aws.String("pBudgetAmount"),
			ParameterValue: aws.String("100"),
		},
	}
}

func createTags(event types.Event) []cloudformationTypes.Tag {
	return []cloudformationTypes.Tag{
		{
			Key:   aws.String("Account"),
			Value: &event.AccountName,
		},
		{
			Key:   aws.String("Domain"),
			Value: &event.Domain,
		},
		{
			Key:   aws.String("Environment"),
			Value: &event.Environment,
		},
		{
			Key:   aws.String("Owner"),
			Value: &event.Owner,
		},
	}
}
