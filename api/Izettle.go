package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// GetAuthkey gets a new auth key from izettle.
func GetAuthkey(username string, password string) string {
	tokenURL := "https://oauth.izettle.com/token"

	clientSecret := os.Getenv("CLIENT_SECRET")
	clientID := os.Getenv("CLIENT_ID")

	log.Println(clientID)
	log.Println(clientSecret)

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
