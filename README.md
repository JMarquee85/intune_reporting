# intune_reporting

This application is intended to provide enrollment information for Intune devices.

We plan to use the Microsoft Export API to get the data, parse it and display it in a simple web interface. GraphAPI calls are also an option if the Export API does not provide the necessary data or to embellish the data.

The web portion is written in Go.

Gorilla mux is being used to handle HTTP requests:
https://github.com/gorilla/mux

### Running the application locally

Ensure Go is installed on your system.

You can create a local .env file or create envirionment variables in your local system.

The following values are required:

CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
TENANT_ID=your_tenant_id

Run `go run .` from the root directory to start the server.
