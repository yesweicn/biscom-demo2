package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SessionResponse represents the structure of the JSON response from the /baseUrl/Session endpoint
type SessionResponse struct {
	Settings          Settings         `json:"settings"`
	TokenInfo         TokenInfo        `json:"tokenInfo"`
}

// Settings represents the "settings" part of the JSON response
type Settings struct {
	WorkflowId                int      `json:"workflowId"`
	RoleId                    int      `json:"roleId"`
	CanDelete                 bool     `json:"canDelete"`
	DepartmentalRole          string   `json:"departmentalRole"`
	MaxPriority               string   `json:"maxPriority"`
}


// TokenInfo represents the "tokenInfo" part of the JSON response
type TokenInfo struct {
	Token              string    `json:"token"`
	ExpirationDate     time.Time `json:"expirationDate"`
	ServerTime         time.Time `json:"serverTime"`
	ServerTimeOffsetMs int       `json:"serverTimeOffsetMs"`
}

func Authenticate(apiURL, username, password string) string {
	// Create a new HTTP client with basic authentication
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// Set basic authentication
	req.SetBasicAuth(username, password)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Check if the response status code is successful (2xx)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error: Non-successful status code received:", resp.Status)
		return ""
	}

	// Decode the JSON response
	var sessionResponse SessionResponse
	err = json.NewDecoder(resp.Body).Decode(&sessionResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	// Now you can use the data retrieved from the API
	fmt.Printf("%+v\n", sessionResponse)
	return sessionResponse.TokenInfo.Token
}
