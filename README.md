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

Clone this repository and navigate to the root directory.

Run `go mod tidy` to install the necessary dependencies.

Run `go run .` from the root directory to start the server.

### Documentation and Resources

[Deploying to Github Actions](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go)

### Example of Output of Endpoints:

`https://graph.microsoft.com/v1.0/deviceManagement/managedDevices`

### Resources and Links

https://github.com/microsoftgraph/msgraph-sdk-go

In order to enable use of the SDK, Allow public client flows must be set to yes in the App Registration.

https://gist.github.com/nikhita/432436d570b89cab172dcf2894465753

https://learn.microsoft.com/en-us/graph/sdks/create-requests?tabs=go

https://learn.microsoft.com/en-us/graph/api/overview?view=graph-rest-1.0

https://github.com/go-echarts/examples

https://github.com/go-echarts/go-echarts

https://go-echarts.github.io/go-echarts/#/
