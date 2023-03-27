package dvc_oauth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Auth0 struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetAuthToken(clientId, clientSecret string) (Auth0, error) {
	url := "https://auth.devcycle.com/oauth/token"

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientSecret + "&audience=https://api.devcycle.com/")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var auth Auth0
	err := json.Unmarshal(body, &auth)
	if err != nil {
		return auth, err
	}

	return auth, nil
}
