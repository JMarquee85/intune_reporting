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

This is a Windows365 Cloud PC.

{
"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#deviceManagement/managedDevices",
"@odata.count": 1,
"@odata.nextLink": "https://graph.microsoft.com/v1.0/deviceManagement/managedDevices?$top=1&$skiptoken=LastDeviceName%3d%27CPC-achak-A17OZ%27%2cLastDeviceId%3d%279b5ca08f-2311-4164-b725-2d29609b35f2%27",
"value": [
{
"id": "9b5ca08f-2311-4164-b725-2d29609b35f2",
"userId": "",
"deviceName": "CPC-achak-A17OZ",
"managedDeviceOwnerType": "company",
"enrolledDateTime": "2023-12-14T18:05:18Z",
"lastSyncDateTime": "2024-04-13T20:23:18Z",
"operatingSystem": "Windows",
"complianceState": "noncompliant",
"jailBroken": "Unknown",
"managementAgent": "mdm",
"osVersion": "10.0.19044.3693",
"easActivated": false,
"easDeviceId": "",
"easActivationDateTime": "0001-01-01T00:00:00Z",
"azureADRegistered": true,
"deviceEnrollmentType": "windowsAzureADJoin",
"activationLockBypassCode": null,
"emailAddress": "",
"azureADDeviceId": "21fa9ef6-f319-4ee2-8c87-097e34eeada7",
"deviceRegistrationState": "registered",
"deviceCategoryDisplayName": "Unknown",
"isSupervised": false,
"exchangeLastSuccessfulSyncDateTime": "0001-01-01T00:00:00Z",
"exchangeAccessState": "none",
"exchangeAccessStateReason": "none",
"remoteAssistanceSessionUrl": null,
"remoteAssistanceSessionErrorDetails": null,
"isEncrypted": false,
"userPrincipalName": "",
"model": "Cloud PC Enterprise 2vCPU/4GB/128GB",
"manufacturer": "Microsoft Corporation",
"imei": "",
"complianceGracePeriodExpirationDateTime": "2023-12-28T22:01:08Z",
"serialNumber": "0000-0014-0657-8375-4585-0309-11",
"phoneNumber": "",
"androidSecurityPatchLevel": "",
"userDisplayName": "",
"configurationManagerClientEnabledFeatures": null,
"wiFiMacAddress": "",
"deviceHealthAttestationState": null,
"subscriberCarrier": "",
"meid": "",
"totalStorageSpaceInBytes": 136844410880,
"freeStorageSpaceInBytes": 110705508352,
"managedDeviceName": "21fa9ef6-f319-4ee2-8c87-097e34eeada7_Windows_12/14/2023_6:05 PM",
"partnerReportedThreatState": "unknown",
"requireUserEnrollmentApproval": null,
"managementCertificateExpirationDate": "2024-12-11T06:59:05Z",
"iccid": null,
"udid": null,
"notes": null,
"ethernetMacAddress": null,
"physicalMemoryInBytes": 0,
"enrollmentProfileName": null,
"deviceActionResults": []
}
]
}
