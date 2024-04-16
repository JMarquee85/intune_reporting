package main

// DeviceInfo Wrapper
// The JSON response information is nested inside a value field
type DeviceInfoWrapper struct {
	Value    []DeviceInfo `json:"value"`
	NextLink string       `json:"@odata.nextLink"`
}

// Struct for devices returned from homeHandler
type DeviceInfo struct {
	ID                     string `json:"id"`
	DeviceId               string `json:"deviceId"`
	DeviceDisplayName      string `json:"displayName"`
	EnrollmentProfileName  string `json:"enrollmentProfileName"`
	EnrollmentType         string `json:"enrollmentType"`
	IsCompliant            bool   `json:"isCompliant"`
	IsManaged              bool   `json:"isManaged"`
	Manufacturer           string `json:"manufacturer"`
	MDMAppId               string `json:"mdmAppId"`
	Model                  string `json:"model"`
	ProfileType            string `json:"profileType"`
	RegistrationDateTime   string `json:"registrationDateTime"`
	LastSignInDateTime     string `json:"approximateLastSignInDateTime"`
	UserID                 string `json:"userId"`
	DeviceName             string `json:"deviceName"`
	DeviceCategory         string `json:"deviceCategory"`
	OperatingSystem        string `json:"operatingSystem"`
	OperatingSystemVersion string `json:"operatingSystemVersion"`
	CreatedDateTime        string `json:"createdDateTime"`
	EnrolledDateTime       string `json:"enrolledDateTime"`
}

type Response struct {
	Value []DeviceInfo `json:"value"`
}

type AzureTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}
