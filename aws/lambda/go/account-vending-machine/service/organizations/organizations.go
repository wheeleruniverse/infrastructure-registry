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
	createAccountOutput := createAccount(event, role)

	environmentOuId := findOrganizationalUnit(&event.OuRootId, &event.Environment)

	targetOuId := findOrganizationalUnit(environmentOuId, &event.OuName)

	moveAccount(createAccountOutput.AccountId, &event.OuRootId, targetOuId)

	return createAccountOutput
}

func createAccount(event types.Event, role *string) types.Account {
	createAccountInput := organizations.CreateAccountInput{
		AccountName:            &event.AccountName,
		Email:                  &event.Owner,
		IamUserAccessToBilling: organizationsTypes.IAMUserAccessToBillingDeny,
		RoleName:               role,
		Tags:                   createTags(event),
	}
	createAccountOutput, createAccountErr := client.CreateAccount(ctx, &createAccountInput)
	if createAccountErr != nil {
		log.Fatalf("organizations.CreateAccount failed because %v", createAccountErr)
	}
	return types.Account{
		AccountId: createAccountOutput.CreateAccountStatus.AccountId,
	}
}

func createTags(event types.Event) []organizationsTypes.Tag {
	return []organizationsTypes.Tag{
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

func findOrganizationalUnit(ouParentId *string, ouTargetName *string) *string {
	listOrganizationalUnitsForParentOutput := listOrganizationalUnitsForParent(ouParentId)

	var ouId *string
	for _, ou := range listOrganizationalUnitsForParentOutput.OrganizationalUnits {
		if ouTargetName == ou.Name {
			ouId = ou.Id
			break
		}
	}
	if ouId == nil {
		log.Fatalf("Could not find 'ouId' for '%v' with 'ouParentId': %v", ouTargetName, ouParentId)
	}
	return ouId
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
