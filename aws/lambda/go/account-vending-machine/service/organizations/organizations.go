package organizations

import (
	"account-vending-machine/types"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationsTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"log"
)

var client *organizations.Client
var ctx = context.TODO()

func Configure(cfg aws.Config) {
	client = organizations.NewFromConfig(cfg)
}

func CreateAccount(event types.Event, role *string) types.Account {
	if client == nil {
		log.Fatalf("organizations.Client is nil")
	}
	createAccountOutput := createAccount(&event.AccountName, &event.Domain, &event.Environment, role)

	listOrganizationalUnitsForParentOutput := listOrganizationalUnitsForParent(&event.OuRootId)

	var ouId *string
	for _, ou := range listOrganizationalUnitsForParentOutput.OrganizationalUnits {
		if event.OuName == *ou.Name {
			ouId = ou.Id
			break
		}
	}
	if ouId == nil {
		log.Fatalf("Could not find 'ouId' for '%v' with 'ouRootId': %v", event.OuName, event.OuRootId)
	}

	moveAccount(createAccountOutput.AccountId, &event.OuRootId, ouId)

	return createAccountOutput
}

func createAccount(
	accountName *string,
	domain *string,
	environment *string,
	role *string,
) types.Account {
	owner := aws.String((*accountName) + "@wheelerswebservices.com")
	createAccountInput := organizations.CreateAccountInput{
		AccountName:            accountName,
		Email:                  owner,
		IamUserAccessToBilling: organizationsTypes.IAMUserAccessToBillingDeny,
		RoleName:               role,
		Tags:                   createTags(accountName, domain, environment, owner),
	}
	createAccountOutput, createAccountErr := client.CreateAccount(ctx, &createAccountInput)
	if createAccountErr != nil {
		log.Fatalf("organizations.CreateAccount failed because %v", createAccountErr)
	}
	return types.Account{
		AccountId: createAccountOutput.CreateAccountStatus.AccountId,
	}
}

func createTags(
	account *string,
	domain *string,
	environment *string,
	owner *string,
) []organizationsTypes.Tag {
	return []organizationsTypes.Tag{
		{
			Key:   aws.String("Account"),
			Value: account,
		},
		{
			Key:   aws.String("Domain"),
			Value: domain,
		},
		{
			Key:   aws.String("Environment"),
			Value: environment,
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

func moveAccount(
	accountId *string,
	destinationParentId *string,
	sourceParentId *string,
) {
	_, moveAccountErr := client.MoveAccount(ctx, &organizations.MoveAccountInput{
		AccountId:           accountId,
		DestinationParentId: destinationParentId,
		SourceParentId:      sourceParentId,
	})
	if moveAccountErr != nil {
		log.Fatalf("organizations.MoveAccount failed because %v", moveAccountErr)
	}
}
