# intune_reporting

This application is intended to provide enrollment information for Intune devices.

We plan to use the Microsoft Export API to get the data, parse it and display it in a simple web interface.

Go is being used here as the backend language.

We are looking into Gorilla Mux or Chi to handle routing for the http requests.

https://github.com/gorilla/mux

### Running the application locally

Ensure Go is installed on your system.

Run `go run main.go` to start the server.
