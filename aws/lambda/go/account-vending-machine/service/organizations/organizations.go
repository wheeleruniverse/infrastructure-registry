package organizations

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"log"
)

var client *organizations.Client
var ctx = context.TODO()

type CreateAccountInput struct {
	AccountName *string
	Domain      *string
	Environment *string
	Role        *string
	OuRootId    *string
	OuName      *string

	accountId *string
	ouId      *string
}

type CreateAccountOutput struct {
	AccountId *string
}

func Configure(cfg aws.Config) {
	client = organizations.NewFromConfig(cfg)
}

func CreateAccount(input CreateAccountInput) CreateAccountOutput {
	if client == nil {
		log.Fatalf("organizations.Client is nil")
	}
	createAccountOutput := createAccount(input)
	input.accountId = createAccountOutput.AccountId

	listOrganizationalUnitsForParentOutput := listOrganizationalUnitsForParent(input.OuRootId)

	var ouId *string
	for _, ou := range listOrganizationalUnitsForParentOutput.OrganizationalUnits {
		if input.OuName == ou.Name {
			ouId = ou.Id
			break
		}
	}
	if ouId == nil {
		log.Fatalf("Could not find 'ouId' for '%v' with 'ouRootId': %v", input.OuName, input.OuRootId)
	}
	input.ouId = ouId

	moveAccount(input)

	return createAccountOutput
}

func createAccount(input CreateAccountInput) CreateAccountOutput {
	accountName := *input.AccountName
	email := aws.String(accountName + "@wheelerswebservices.com")
	createAccountInput := organizations.CreateAccountInput{
		AccountName:            input.AccountName,
		Email:                  email,
		IamUserAccessToBilling: types.IAMUserAccessToBillingDeny,
		RoleName:               input.Role,
		Tags:                   createTags(input.AccountName, input.Environment, input.Domain, email),
	}
	createAccountOutput, createAccountErr := client.CreateAccount(ctx, &createAccountInput)
	if createAccountErr != nil {
		log.Fatalf("organizations.CreateAccount failed because %v", createAccountErr)
	}
	return CreateAccountOutput{
		AccountId: createAccountOutput.CreateAccountStatus.AccountId,
	}
}

func createTags(
	account *string,
	environment *string,
	domain *string,
	owner *string,
) []types.Tag {
	return []types.Tag{
		{
			Key:   aws.String("Account"),
			Value: account,
		},
		{
			Key:   aws.String("Environment"),
			Value: environment,
		},
		{
			Key:   aws.String("Domain"),
			Value: domain,
		},
		{
			Key:   aws.String("Owner"),
			Value: owner,
		},
	}
}

func listOrganizationalUnitsForParent(ouRootId *string) *organizations.ListOrganizationalUnitsForParentOutput {
	listOrganizationalUnitsForParentInput := organizations.ListOrganizationalUnitsForParentInput{
		ParentId: ouRootId,
	}
	listOrganizationalUnitsForParentOutput, listOrganizationalUnitsForParentErr := client.ListOrganizationalUnitsForParent(
		ctx, &listOrganizationalUnitsForParentInput,
	)
	if listOrganizationalUnitsForParentErr != nil {
		log.Fatalf("organizations.ListOrganizationalUnitsForParent failed because %v", listOrganizationalUnitsForParentErr)
	}
	return listOrganizationalUnitsForParentOutput
}

func moveAccount(input CreateAccountInput) {
	_, moveAccountErr := client.MoveAccount(ctx, &organizations.MoveAccountInput{
		AccountId:           input.accountId,
		DestinationParentId: input.OuRootId,
		SourceParentId:      input.ouId,
	})
	if moveAccountErr != nil {
		log.Fatalf("organizations.MoveAccount failed because %v", moveAccountErr)
	}
}
