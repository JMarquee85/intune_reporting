package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/joho/godotenv"
)

var clientID string
var tenantID string
var clientSecret string

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientID = os.Getenv("CLIENT_ID")
	tenantID = os.Getenv("TENANT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
}

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

	// w.Write([]byte("This is the device test handler!"))

	// Get All Devices
	var allDevices []DeviceInfo

	apiURL := "https://graph.microsoft.com/v1.0/deviceManagement/managedDevices?$top=1"
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

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a test handler!"))
}

func enrollmentsHandler(w http.ResponseWriter, r *http.Request) {
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
	// Temporarily changed to 300 days for testing purposes
	sevenDaysAgo := time.Now().Add(-300 * 24 * time.Hour)
	var androidEnrollDates []time.Time
	var iOSEnrollDates []time.Time

	for _, device := range allDevices {
		enrolledDateTime, err := time.Parse(time.RFC3339, device.EnrolledDateTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if enrolledDateTime.After(sevenDaysAgo) {
			if device.OperatingSystem == "Android" {
				androidEnrollDates = append(androidEnrollDates, enrolledDateTime)
			} else if device.OperatingSystem == "iOS" {
				iOSEnrollDates = append(iOSEnrollDates, enrolledDateTime)
			}
		}
	}

	androidEnrollmentCount := len(androidEnrollDates)
	iOSEnrollmentCount := len(iOSEnrollDates)
	log.Printf("Found %d Android enrollments in the last seven days", androidEnrollmentCount)
	log.Printf("Found %d iOS enrollments in the last seven days", iOSEnrollmentCount)

	w.Write([]byte(fmt.Sprintf("Found %d Android and %d iOS enrollments in the last seven days", androidEnrollmentCount, iOSEnrollmentCount)))

	// Add a couple line breaks
	w.Write([]byte("<br><br>"))

	// Attempting a line graph demo
	// Create a new line instance
	line := charts.NewLine()

	// Prepare the data
	androidData := make([]opts.LineData, 0)
	iOSData := make([]opts.LineData, 0)
	xAxisLabels := make([]string, 0)

	// Calculate the number of weeks to display
	numWeeks := 10 // change this to the number of weeks you want to display

	for i := 0; i < numWeeks; i++ {
		// Calculate the start and end of the week
		startOfWeek := time.Now().AddDate(0, 0, -7*i).Format("2006-01-02")
		endOfWeek := time.Now().AddDate(0, 0, -7*(i+1)).Format("2006-01-02")

		// Add the week to the X-axis labels
		xAxisLabels = append([]string{startOfWeek + " to " + endOfWeek}, xAxisLabels...)

		// Count the number of enrollments for the week
		androidCount := 0
		iOSCount := 0
		for _, date := range androidEnrollDates {
			if date.After(time.Now().AddDate(0, 0, -7*(i+1))) && date.Before(time.Now().AddDate(0, 0, -7*i)) {
				androidCount++
			}
		}
		for _, date := range iOSEnrollDates {
			if date.After(time.Now().AddDate(0, 0, -7*(i+1))) && date.Before(time.Now().AddDate(0, 0, -7*i)) {
				iOSCount++
			}
		}

		// Add the counts to the data
		androidData = append([]opts.LineData{{Value: androidCount}}, androidData...)
		iOSData = append([]opts.LineData{{Value: iOSCount}}, iOSData...)
	}

	// Set the options
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Enrollment Trends",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "value",
		}),
	)

	// Add the data
	line.SetXAxis(xAxisLabels).
		AddSeries("Android", androidData).
		AddSeries("iOS", iOSData)

	// Render the chart
	page := components.NewPage()
	page.AddCharts(line)
	err = page.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
