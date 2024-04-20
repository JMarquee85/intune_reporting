package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func workspaceOneAuth(w http.ResponseWriter, r *http.Request) (string, error) {

	// Help URL
	// https://as1352.awmdm.com/api/help
	// https://docs.vmware.com/en/VMware-Workspace-ONE-UEM/services/System_Settings_On_Prem/GUID-AWT-SYSTEM-ADVANCED-API-REST.html
	// https://as1352.awmdm.com/api/system/info

	// Test logging
	// fmt.Printf("Workspace One URL: %v\n", workspaceOneUrl)
	// fmt.Printf("Workspace One API Key: %v\n", workspaceOneApiKey)

	req, err := http.NewRequest("GET", workspaceOneUrl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	req.Header.Set("id", workspaceOneApiKey)

	// req.Header.Set("Authorization", "Bearer "+workspaceOneApiKey)
	req.Header.Set("Accept", "application/json")
	// req.Header.Set("Content-Length", "0")
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Host", "as1352.awmdm.com")
	req.Header.Set("aw-tenant-code", workspaceOneApiKey)

	// Print the request
	// fmt.Printf("Request: %v\n", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error response: %v\n", resp)
		http.Error(w, "Workspace One authentication failed", http.StatusUnauthorized)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	message := fmt.Sprintf("Workspace One authentication successful: %s", body)
	// fmt.Fprint(w, message)
	return message, nil

}
