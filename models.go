package main

// DeviceInfo Wrapper
// The JSON response information is nested inside a value field
type DeviceInfoWrapper struct {
	Value    []DeviceInfo `json:"value"`
	NextLink string       `json:"@odata.nextLink"`
}

// Struct for devices returned from homeHandler
type DeviceInfo struct {
	UserID           string `json:"userId"`
	DeviceName       string `json:"deviceName"`
	OperatingSystem  string `json:"operatingSystem"`
	EnrolledDateTime string `json:"enrolledDateTime"`
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
