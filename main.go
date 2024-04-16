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
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/devices", deviceTest).Methods("GET")
	r.HandleFunc("/test", testHandler).Methods("GET")
	r.HandleFunc("/enrollments", enrollmentsHandler).Methods("GET")

	fmt.Println("Server is starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
