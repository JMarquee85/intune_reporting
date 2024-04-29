package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type WorkspaceOneTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func getWorkspaceOneToken(clientID string, clientSecret string, workspaceOneTokenUrl string) (string, time.Time, error) {

	// prepare the form data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// make the request
	resp, err := http.PostForm(workspaceOneTokenUrl, data)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", time.Time{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP response status: %s", resp.Status)
		log.Printf("HTTP response body: %s", string(body))
		return "", time.Time{}, err
	}

	var tokenResponse WorkspaceOneTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return "", time.Time{}, err
	}

	expiryTime := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return tokenResponse.AccessToken, expiryTime, nil

}
