package handlers

import (
	"net/http"
)

// HomeHandler is the handler for the home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to the home page!"))
}
