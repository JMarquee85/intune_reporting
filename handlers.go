package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

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
	androidEnrollDates, iOSEnrollDates, err := filterEnrollments(allDevices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create graph of enrollments
	renderIntuneEnrollmentGraph(w, androidEnrollDates, iOSEnrollDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func migrationAssistantHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/template.html"))

	// Add high device check. Offer to remove oldest devices if over 10.

	data := PageData{
		Title:  "Marriott Intune Mobile Migration Assistant",
		Header: "Marriott Intune Mobile Migration Assistant",
		Content: template.HTML(`
        This tool will help you migrate your mobile devices from Workspace ONE to Microsoft Intune.<br><br>
        <form id="eidForm" action="/migrationp1" method="post">
            Please enter your EID:<br>
            <input type="text" id="eidInput" name="EID"><br>
            <button type="button" id="submitBtn">Submit</button>
            <div id="loadingIcon" style="display: none;">Loading...</div>
        </form>
        <script>
	document.getElementById('submitBtn').addEventListener('click', function(event) {
    event.preventDefault();

    document.getElementById('loadingIcon').style.display = 'block';

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/migrationp1', true);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.onload = function() {
        if (this.status == 200) {
            var response = JSON.parse(this.responseText);
            console.log(response.message);
            if (response.success) {
                window.location.href = '/migrationp2';
            } else {
                // Display error message
            }
        }
        document.getElementById('loadingIcon').style.display = 'none';
    };
    xhr.send('EID=' + encodeURIComponent(document.getElementById('eidInput').value));
});
</script>
        `),
	}

	tmpl.Execute(w, data)
}

func migrationP1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	eid := r.FormValue("EID")

	// Here we will look for registered iOS or Android devices in Workspace ONE
	// Instead of printing the message below, we will do the stuff.

	success := false // Set this to true or false based on whether the operation was successful

	response := map[string]interface{}{
		"message": fmt.Sprintf("Will check here for registrations in WorkspaceOne related to user: %s", eid),
		"success": success,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func migrationP2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/template.html"))

	data := PageData{
		Title:  "Marriott Intune Mobile Migration Assistant",
		Header: "Marriott Intune Mobile Migration Assistant",
		Content: template.HTML(`
		<p>Migration Assistant Step 2</p>
		`),
	}

	tmpl.Execute(w, data)
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
