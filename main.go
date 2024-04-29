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
var workspaceOneClientID string
var workspaceOneClientSecret string
var workspaceOneTokenUrl string

// var workspaceOneUrl string

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientID = os.Getenv("CLIENT_ID")
	tenantID = os.Getenv("TENANT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	workspaceOneTokenUrl = os.Getenv("WORKSPACE_ONE_TOKEN_URL")
	workspaceOneClientID = os.Getenv("WORKSPACE_ONE_CLIENT_ID")
	workspaceOneClientSecret = os.Getenv("WORKSPACE_ONE_CLIENT_SECRET")
	// workspaceOneUrl = os.Getenv("WORKSPACE_ONE_URL")

}

func main() {

	// Create a new instance of mux router
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/devices", deviceTest).Methods("GET")
	r.HandleFunc("/reports", reportingHandler).Methods("GET")
	r.HandleFunc("/workspaceonefailed", workspaceOneFailedHandler).Methods("GET")

	// Workspace One API Testing
	r.HandleFunc("/workspaceone", workspaceOneHandler).Methods("GET")

	fmt.Println("Server is starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
