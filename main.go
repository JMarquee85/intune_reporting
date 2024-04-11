package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new instance of mux router
	r := mux.NewRouter()

	// Define the route(s)
	r.HandleFunc("/", homeHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
