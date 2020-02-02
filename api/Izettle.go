package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Constants

const tokenURL = "https://oauth.izettle.com/token"
const transactionURL = "https://finance.izettle.com/organizations/us/accounts/liquid/transactions?start=%[1]s&end=%[2]s"
const purchaseURL = "https://purchase.izettle.com/purchases/v2?startDate=%[1]s&endDate=%[2]s&descending=true"

// Interval of days to get data from izettle, depends on how often payouts are made.
const dateInterval = 14

// GetAuthkey gets a new auth key from izettle.
func GetAuthkey(username string, password string) (string, error) {

	clientSecret := os.Getenv("CLIENT_SECRET")
	clientID := os.Getenv("CLIENT_ID")

	formData := url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password},
	}

	resp, err := http.PostForm(tokenURL, formData)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		var error string
		json.NewDecoder(resp.Body).Decode(&error)
		return "", errors.New(error)
	}

	var result tokenResponse

	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken, nil
}

func izettleGetRequest(startDate string, endDate string, auth string, URL string) (*http.Response, error) {

	formattedURL := fmt.Sprintf(URL, startDate, endDate)

	client := http.Client{}
	req, err := http.NewRequest("GET", formattedURL, nil)
	if err != nil {
		return nil, err
	}

	// Use auth token
	req.Header.Set("Authorization", "Bearer "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

// GetTransactions returns the transactions made from (endDate - 2 weeks) until endDate
func GetTransactions(startDate string, endDate string, auth string) ([]Transaction, error) {

	resp, err := izettleGetRequest(startDate, endDate, auth, transactionURL)
	if err != nil {
		return nil, err
	}

	var result transactionData
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Data, nil

}

// GetPurchases returns the purchases made from (endDate - 2 weeks) until endDate
// I give up
func GetPurchases(startDate string, endDate string, auth string) ([]Purchase, error) {

	resp, err := izettleGetRequest(startDate, endDate, auth, purchaseURL)
	if err != nil {
		return nil, err
	}

	var result purchaseData
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Purchases, nil

}
