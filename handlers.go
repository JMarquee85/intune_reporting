package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
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

// Checks for valid WorkspaceOne token and gets another if expired
// This and the above function could be combined into one with some kind of argument passed in to determine which token to get
// Maybe take Azure or WorkspaceOne as an argument and add an if statement in the function to determine which token to get
func ensureWorkspaceOneToken() error {
	// print expiry time
	fmt.Printf("Current time: %s\n", time.Now())
	fmt.Printf("Token expiry time: %s\n", expiryTime)

	if time.Now().After(expiryTime) {
		var err error
		accessToken, expiryTime, err = getWorkspaceOneToken(workspaceOneClientID, workspaceOneClientSecret, workspaceOneTokenUrl)
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

func deviceTest(w http.ResponseWriter, r *http.Request) {
	err := ensureAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Get All Devices
	var allDevices []IntuneDeviceInfo

	apiURL := "https://graph.microsoft.com/v1.0/deviceManagement/managedDevices"
	// This should handle pagination
	for apiURL != "" {
		body, err := makeGraphAPIRequest(accessToken, apiURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var intuneDeviceInfoWrapper IntuneDeviceInfoWrapper
		err = json.Unmarshal(body, &intuneDeviceInfoWrapper)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allDevices = append(allDevices, intuneDeviceInfoWrapper.Value...)

		apiURL = intuneDeviceInfoWrapper.NextLink
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
	var allDevices []IntuneDeviceInfo

	apiURL := "https://graph.microsoft.com/v1.0/deviceManagement/managedDevices"
	for apiURL != "" {
		body, err := makeGraphAPIRequest(accessToken, apiURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var deviceInfoWrapper IntuneDeviceInfoWrapper
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
	err = renderBarChart(w, "All Regions", intuneData, workspaceOneData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Stacked Bar AMER
	err = renderBarChart(w, "AMER", intuneData, workspaceOneData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Stacked Bar EMEA
	err = renderBarChart(w, "EMEA", intuneData, workspaceOneData)
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

	// Set header as HTML
	// w.Header().Set("Content-Type", "text/html")

	// Check for valid WorkspaceOne Token
	err := ensureWorkspaceOneToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Make a call to the WorkspaceOne API
	apiURL := workspaceOneUrl + "/API/mdm/devices/search"
	var allDevices []WorkspaceOneDeviceInfo
	pageSize := 500
	curPage := 0

	for {
		// Make a call to the WorkspaceOne API
		body, err := makeWorkspaceOneRequest(accessToken, fmt.Sprintf("%s?page=%d&pagesize=%d", apiURL, curPage, pageSize))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the JSON response into a Response value
		var resp WorkspaceOneResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allDevices = append(allDevices, resp.Devices...)

		// Check if there are more results to fetch
		finalPage := int(math.Ceil(float64(resp.Total)/float64(pageSize))) - 1
		if curPage >= finalPage {
			break
		}

		// Increment the current page
		curPage++
	}

	deviceCount := len(allDevices)

	log.Printf("Found %d devices", deviceCount)

	// Convert allDevices into JSON
	jsonResponse, err := json.MarshalIndent(allDevices, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write jsonResponse to the HTTP response
	w.Write(jsonResponse)
}
