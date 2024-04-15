package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var clientID string
var tenantID string
var clientSecret string

func init() {
	clientID = os.Getenv("CLIENT_ID")
	tenantID = os.Getenv("TENANT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
}

var (
	accessToken string
	expiryTime  time.Time
)

func ensureAccessToken() error {
	if time.Now().After(expiryTime) {
		var err error
		accessToken, expiryTime, err = getAzureToken(clientID, clientSecret, tenantID)
		if err != nil {
			return err
		}
	}
	return nil

}

// HomeHandler is the handler for the home route
func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to the home page!"))
}

func deviceTest(w http.ResponseWriter, r *http.Request) {
	err := ensureAccessToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get All Devices
	var allDevices []DeviceInfo

	apiURL := "https://graph.microsoft.com/v1.0/deviceManagement/managedDevices"
	for apiURL != "" {
		body, err := makeGraphAPIRequest(accessToken, apiURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var deviceInfoWrapper DeviceInfoWrapper
		err = json.Unmarshal(body, &deviceInfoWrapper)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allDevices = append(allDevices, deviceInfoWrapper.Value...)

		apiURL = deviceInfoWrapper.NextLink
	}

	var devicesToSend []DeviceInfo
	if len(allDevices) > 25 {
		devicesToSend = allDevices[:25]
	} else {
		devicesToSend = allDevices
	}

	jsonResponse, err := json.MarshalIndent(devicesToSend, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a test handler!"))
}
