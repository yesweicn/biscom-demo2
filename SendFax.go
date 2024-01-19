package BiscomLibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"time"
)

// Define any additional structs needed for SendFax
// Attachment represents an attachment in the JSON structure.
type Attachment struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Recipient represents a recipient in the JSON structure.
type Recipient struct {
	FaxNumber   string `json:"faxNumber"`
	DeliveryType string `json:"deliveryType"`
}

// YourStruct represents the overall structure of your JSON data.
type FaxData struct {
	Attachments  []Attachment `json:"attachments"`
	Subject      string       `json:"subject"`
	CoverPage    string       `json:"coverPage"`
	Memo         string       `json:"memo"`
	DeliveryTime time.Time    `json:"deliveryTime"`
	Priority     string       `json:"priority"`
	Recipients   []Recipient  `json:"recipients"`
}

func SendFax(apiURL, token, payload string) error {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request (example assumes a POST request)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set the Content-Type header to indicate that the request body contains JSON data
	req.Header.Set("Content-Type", "application/json")
	// Set the Authorization header with the Bearer Token
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code is successful (2xx)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error: Non-successful status code received:", resp.Status)
		return fmt.Errorf("Error: %s", resp.Status);
	}

	// Handle the response as needed
	fmt.Println("Request successful!")
	return nil
}
