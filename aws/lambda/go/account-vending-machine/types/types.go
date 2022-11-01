package types

type Account struct {
	AccountId *string
}

type Event struct {
	AccountName string
	Domain      string
	Environment string
	OuRootId    string
	OuName      string
}
