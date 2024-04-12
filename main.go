package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var clientID string
var tenantID string
var clientSecret string

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID = os.Getenv("CLIENT_ID")
	tenantID = os.Getenv("TENANT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	// Create a new instance of mux router
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")

	fmt.Println("Server is starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	accessToken, err := getAzureToken(clientID, clientSecret, tenantID)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Use accessToken to make requests	to Microsoft Graph API
	// fmt.Fprintf(w, "Client ID: %s\n", clientID)
	// fmt.Fprintf(w, "Tenant ID: %s\n", tenantID)
	// fmt.Fprintf(w, "Client Secret: %s\n", clientSecret)
	fmt.Fprintf(w, "Access Token: %s\n", accessToken)
	fmt.Println(accessToken)

}
