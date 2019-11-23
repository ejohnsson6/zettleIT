package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

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

// TransactionData is a wrapper for Transactions
type TransactionData struct {
	Data []Transaction `json:"data"`
}

// Constants

const tokenURL = "https://oauth.izettle.com/token"
const transactionURL = "https://finance.izettle.com/organizations/us/accounts/preliminary/transactions?start=%[1]s&end=%[2]s"

// GetAuthkey gets a new auth key from izettle.
func GetAuthkey(username string, password string) string {

	clientSecret := os.Getenv("CLIENT_SECRET")
	clientID := os.Getenv("CLIENT_ID")

	log.Println("Read Client ID from Environment Variables: " + clientID)
	log.Println("Read Client Secret from Environment Variables: " + clientSecret)

	formData := url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password},
	}

	resp, err := http.PostForm(tokenURL, formData)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		var error map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&error)
		log.Fatalln(error)
	}

	var result tokenResponse

	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken
}

// GetTransactions returns the transactions made from (endDate - 2 weeks) until endDate
func GetTransactions(startDate string, endDate string, auth string) []Transaction {

	formattedURL := fmt.Sprintf(transactionURL, startDate, endDate)

	println(formattedURL)

	client := http.Client{}
	req, err := http.NewRequest("GET", formattedURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Use auth token
	req.Header.Set("Authorization", "Bearer "+auth)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	var result TransactionData
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Data

}
