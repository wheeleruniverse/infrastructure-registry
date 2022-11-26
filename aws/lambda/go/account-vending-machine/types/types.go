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
	Owner       string
}

type Request struct {
	LogicalResourceId  string `json:"LogicalResourceId"`
	PhysicalResourceId string `json:"PhysicalResourceId"`
	RequestId          string `json:"RequestId"`
	RequestType        string `json:"RequestType"`
	ResourceType       string `json:"ResourceType"`
	ResponseURL        string `json:"ResponseURL"`
	StackId            string `json:"StackId"`

	ResourceProperties Event
}

type Response struct {
	LogicalResourceId  string `json:"LogicalResourceId"`
	PhysicalResourceId string `json:"PhysicalResourceId"`
	Reason             string `json:"Reason"`
	RequestId          string `json:"RequestId"`
	StackId            string `json:"StackId"`
	Status             string `json:"Status"`

	/*
		use this field to return any response data

		Data struct {}
	*/
}
