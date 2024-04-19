package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	accessToken string
	expiryTime  time.Time
)

func ensureAccessToken() error {
	// print expiry time
	fmt.Printf("Current time: %s\n", time.Now())
	fmt.Printf("Token expiry time: %s\n", expiryTime)

	if time.Now().After(expiryTime) {
		var err error
		accessToken, expiryTime, err = getAzureToken(clientID, clientSecret, tenantID)
		if err != nil {
			log.Printf("Error getting access token: %v", err)
			return err
		}
		log.Printf("Successfully refreshed access token.")
	} else {
		log.Printf("Access token is still valid.")
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
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Get All Devices
	var allDevices []DeviceInfo

	apiURL := "https://graph.microsoft.com/v1.0/deviceManagement/managedDevices"
	// This should handle pagination
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

	jsonResponse, err := json.MarshalIndent(allDevices, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deviceCount := len(allDevices)
	log.Printf("Found %d devices", deviceCount)

	w.Write(jsonResponse)
}

func reportingHandler(w http.ResponseWriter, r *http.Request) {
	// Set header as HTML
	w.Header().Set("Content-Type", "text/html")

	err := ensureAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Get All Enrollments
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

	// Filter for Android and iOS devices enrolled in the last seven days
	androidEnrollDates, iOSEnrollDates, err := filterEnrollments(allDevices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Intune Enrollment Line Graph
	renderIntuneEnrollmentGraph(w, androidEnrollDates, iOSEnrollDates)

	// Feeding this dummy data
	// Will later call functions getUserCountIntuneGroup to get the actual numbers
	intuneData := [3]int{100, 200, 300}
	workspaceOneData := [3]int{150, 250, 400}

	// Stacked Bar All Regions
	err = renderBarChart(w, "All Regions Device Comparison\n", intuneData, workspaceOneData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Stacked Bar AMER
	err = renderBarChart(w, "AMER Device Comparison\n", intuneData, workspaceOneData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Stacked Bar EMEA
	err = renderBarChart(w, "EMEA Device Comparison\n", intuneData, workspaceOneData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func workspaceOneFailedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/template.html"))

	data := PageData{
		Title:  "Marriott Intune Mobile Migration Assistant",
		Header: "Marriott Intune Mobile Migration Assistant",
		Content: template.HTML(`
		The Workspace ONE migration process failed. Please contact your IT department for assistance.
		`),
	}

	tmpl.Execute(w, data)
}

// WorkspaceOne API Handler Testing
func workspaceOneHandler(w http.ResponseWriter, r *http.Request) {
	message, err := workspaceOneAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, message)
}
