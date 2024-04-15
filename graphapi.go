package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetAccessToken(clientID, clientSecret, tenantID string) (string, time.Time, error) {
	accessToken, expiryTime, err := getAzureToken(clientID, clientSecret, tenantID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get access token: %w", err)
	}
	return accessToken, expiryTime, nil
}
func makeGraphAPIRequest(accessToken, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP response status: %s", resp.Status)
		log.Printf("HTTP response body: %s", string(body))
		return nil, fmt.Errorf("HTTP request failed with status %s", resp.Status)
	}

	var errorResponse struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	// Try to unmarshal the response body into the errorResponse struct
	json.Unmarshal(body, &errorResponse)

	if errorResponse.Error.Code == "Authorization_RequestDenied" {
		return nil, fmt.Errorf("insufficient privileges to make this GraphAPI call")
	}

	return body, nil
}
