package main

import (
	"time"
)

func filterEnrollments(allDevices []IntuneDeviceInfo) ([]time.Time, []time.Time, error) {
	sevenDaysAgo := time.Now().Add(7 * 24 * time.Hour)
	var androidEnrollDates []time.Time
	var iOSEnrollDates []time.Time

	for _, device := range allDevices {
		enrolledDateTime, err := time.Parse(time.RFC3339, device.EnrolledDateTime)
		if err != nil {
			return nil, nil, err
		}

		if enrolledDateTime.After(sevenDaysAgo) {
			if device.OperatingSystem == "Android" {
				androidEnrollDates = append(androidEnrollDates, enrolledDateTime)
			} else if device.OperatingSystem == "iOS" {
				iOSEnrollDates = append(iOSEnrollDates, enrolledDateTime)
			}
		}
	}

	return androidEnrollDates, iOSEnrollDates, nil
}

// Stacked Bar info retrieval

// Get all users in the various Intune region groups
// Placeholder functions just to return numbers.

// ////////////////////////////////////////////////
// When polling for actual groups, will refactor this into accepting a group ID and return the counts
// for that group.
// For now, just returning junk data to get a graph going.
// func getUserCountIntuneGroup(accessToken string, groupID string) (int, error) {
// 	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/members", groupID)

// 	// Make the request
// 	body, err := makeGraphAPIRequest(accessToken, url)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Parse the response
// 	var result map[string]interface{}
// 	json.Unmarshal(body, &result)

// 	// Get the members from the response
// 	members, ok := result["value"].([]interface{})
// 	if !ok {
// 		return 0, fmt.Errorf("failed to get members from response")
// 	}
// 	// Return count of members
// 	return len(members), nil
// }

// Write something similar for getting user count to the above once you have WS1 access to do so.
