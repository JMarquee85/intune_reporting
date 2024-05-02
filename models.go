package main

import (
	"html/template"
)

// DeviceInfo Wrapper
// The JSON response information is nested inside a value field
type IntuneDeviceInfoWrapper struct {
	Value    []IntuneDeviceInfo `json:"value"`
	NextLink string             `json:"@odata.nextLink"`
}

// Struct for devices returned from homeHandler
type IntuneDeviceInfo struct {
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

type WorkspaceOneDeviceInfo struct {
	EasIds struct {
		EasId []string `json:"easId"`
	} `json:"EasIds"`
	TimeZone           string `json:"TimeZone"`
	Udid               string `json:"Udid"`
	AssetNumber        string `json:"AssetNumber"`
	DeviceFriendlyName string `json:"DeviceFriendlyName"`
	DeviceReportedName string `json:"DeviceReportedName"`
	LocationGroupId    struct {
		Id struct {
			Value int `json:"value"`
		} `json:"Id"`
		Name string `json:"Name"`
		Uuid string `json:"Uuid"`
	} `json:"LocationGroupId"`
	LocationGroupName string `json:"LocationGroupName"`
	UserName          string `json:"UserName"`
	UserEmailAddress  string `json:"UserEmailAddress"`
	Ownership         string `json:"Ownership"`
	Platform          string `json:"Platform"`
	OperatingSystem   string `json:"OperatingSystem"`
	LastSeen          string `json:"LastSeen"`
	EnrollmentStatus  string `json:"EnrollmentStatus"`
	ComplianceStatus  string `json:"ComplianceStatus"`
	LastEnrolledOn    string `json:"LastEnrolledOn"`
	IsSupervised      bool   `json:"IsSupervised"`
	EnrolledViaDEP    bool   `json:"EnrolledViaDEP"`
}

type Response struct {
	Value []IntuneDeviceInfo `json:"value"`
}

type WorkspaceOneResponse struct {
	Devices []WorkspaceOneDeviceInfo `json:"Devices"`
	Total   int                      `json:"Total"`
}

type AzureTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type PageData struct {
	Title   string
	Header  string
	Content template.HTML
}
