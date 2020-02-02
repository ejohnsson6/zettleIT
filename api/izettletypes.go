package api

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// A Transaction represents an izettle transaction
type Transaction struct {
	Timestamp string `json:"timestamp"`
	Amount    int    `json:"amount"`
	Type      string `json:"originatorTransactionType"`
	UUID      string `json:"originatingTransactionUuid"`
}

type transactionData struct {
	Data []Transaction `json:"data"`
}

type payment struct {
	UUID string `json:"uuid"`
}

// A Purchase represents an izettle purchase
type Purchase struct {
	UserDisplayName string    `json:"userDisplayName"`
	Payments        []payment `json:"payments"`
}

type purchaseData struct {
	Purchases []Purchase `json:"purchases"`
}
