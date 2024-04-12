package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AzureTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

func getAzureToken(clientID string, clientSecret string, tenantID string) (string, error) {
	requestURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(requestURL, data)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP response status: %s", resp.Status)
		log.Printf("HTTP response body: %s", string(body))
		return "", fmt.Errorf("HTTP request failed with status %s", resp.Status)
	}

	var tokenResponse AzureTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
