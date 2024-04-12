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

	// Get the access token
	accessToken, err := getAccessToken(clientID, clientSecret, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Making some test calls to Graph API

	body, err := makeGraphAPIRequest(accessToken, "https://graph.microsoft.com/v1.0/devices")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Print the response body
	fmt.Println(string(body))

}
