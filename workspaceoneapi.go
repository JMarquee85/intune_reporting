package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func makeWorkspaceOneRequest(accessToken, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("insufficient privileges to make this WorkspaceOne API call")
	}

	return body, nil
}
