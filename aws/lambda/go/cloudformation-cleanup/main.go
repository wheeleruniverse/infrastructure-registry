package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	lambda.Start(handleRequest)
}

func getTimeToLive(fallback float64) float64 {
	strVal, success := os.LookupEnv("TTL")
	if success {
		intVal, atoiErr := strconv.Atoi(strVal)
		if atoiErr != nil {
			log.Printf("Atoi failed because %v", atoiErr)
			return fallback
		}
		return float64(intVal)
	}
	return fallback
}

func handleRequest() ([]string, error) {
	cfg, cfgErr := config.LoadDefaultConfig(context.TODO())
	if cfgErr != nil {
		log.Fatalf("LoadDefaultConfig failed because %v", cfgErr)
	}
	cloudformationClient := cloudformation.NewFromConfig(cfg)

	listStacksInput := cloudformation.ListStacksInput{
		StackStatusFilter: []types.StackStatus{types.StackStatusCreateComplete},
	}
	listStacksOutput, listStacksErr := cloudformationClient.ListStacks(context.TODO(), &listStacksInput)
	if listStacksErr != nil {
		log.Fatalf("ListStacks failed because %v", listStacksErr)
	}

	stacksToDelete := make([]string, 0)
	ttl := getTimeToLive(90)
	for _, val := range listStacksOutput.StackSummaries {
		nameMatches := strings.HasPrefix(*val.StackName, "cat-cloud9-")
		timeMatches := time.Now().UTC().Sub(*val.CreationTime).Minutes() > ttl

		if nameMatches && timeMatches {
			stacksToDelete = append(stacksToDelete, *val.StackName)
		}
	}
	for _, val := range stacksToDelete {
		log.Printf("DeleteStack %v", val)
		deleteStackInput := cloudformation.DeleteStackInput{StackName: &val}
		_, deleteStackErr := cloudformationClient.DeleteStack(context.TODO(), &deleteStackInput)

		if deleteStackErr != nil {
			log.Fatalf("DeleteStack failed because %v", deleteStackErr)
		}
	}
	return stacksToDelete, nil
}
