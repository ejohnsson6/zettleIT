package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tkanos/gonfig"
)

type secrets struct {
	clientSecret string `json:Client_Secret`
	clientID     string `json:Client_ID`
}

// GetAuthkey gets a new auth key from izettle.
func GetAuthkey(username string, password string) string {
	url := "https://oauth.izettle.com/token"

	secrets := secrets{}
	err := gonfig.GetConf("./apis/secrets.json", &secrets)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("grant_type", "password")
	req.Header.Add("client_id", secrets.clientID)
	req.Header.Add("client_secret", secrets.clientSecret)
	req.Header.Add("username", username)
	req.Header.Add("password", password)
	resp, _ := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	return "fuckoff"
}
