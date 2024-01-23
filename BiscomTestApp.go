package main

import (
	"fmt"	
	"bytes"
	"time"
	"encoding/json"	
	"encoding/base64"	
	"net/http"
	"io/ioutil"
	"log"
)

// SessionResponse represents the structure of the JSON response from the /baseUrl/Session endpoint
type SessionResponse struct {
	Settings          Settings         `json:"settings"`
	NoteMetadata      []NoteMetadata   `json:"noteMetadata"`
	AvailableFolders  []AvailableFolder `json:"availableFolders"`
	TokenInfo         TokenInfo        `json:"tokenInfo"`
}

// Settings represents the "settings" part of the JSON response
type Settings struct {
	WorkflowId                int      `json:"workflowId"`
	RoleId                    int      `json:"roleId"`
	CanDelete                 bool     `json:"canDelete"`
	DepartmentalRole          string   `json:"departmentalRole"`
	MaxPriority               string   `json:"maxPriority"`
	AvailableDeliveryModes    []string `json:"availableDeliveryModes"`
	CanMoveOutOfArchive       bool     `json:"canMoveOutOfArchive"`
	CanCopyOrMoveToArchive    bool     `json:"canCopyOrMoveToArchive"`
	AutoArchiveTx             bool     `json:"autoArchiveTx"`
	AutoArchiveRx             bool     `json:"autoArchiveRx"`
	CanViewArchive            bool     `json:"canViewArchive"`
	CanViewPublicPhonebook    bool     `json:"canViewPublicPhonebook"`
	CanViewDepartmentPhonebook bool     `json:"canViewDepartmentPhonebook"`
	CanViewPrivatePhonebook   bool     `json:"canViewPrivatePhonebook"`
	DepartmentFolderId        int      `json:"departmentFolderId"`
	UserMailboxId             int      `json:"userMailboxId"`
}

// NoteMetadata represents the "noteMetadata" part of the JSON response
type NoteMetadata struct {
	Id         int `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

// AvailableFolder represents the "availableFolders" part of the JSON response
type AvailableFolder struct {
	Path    string `json:"path"`
	Id      int    `json:"id"`
	JobType string `json:"jobType"`
}


// TokenInfo represents the "tokenInfo" part of the JSON response
type TokenInfo struct {
	Token              string    `json:"token"`
	ExpirationDate     string `json:"expirationDate"`
	ServerTime         string `json:"serverTime"`
	ServerTimeOffsetMs float32       `json:"serverTimeOffsetMs"`
}

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
		fmt.Println("Error: Non-successful status code received:", resp)
		return ""
	} 

	//For troubleshooting - Read and print the response body
	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}
	fmt.Println("Response Body:", body.String())

	// Decode the JSON response
	var sessionResponse SessionResponse
	//cannot user resp.Body, that stream can only be read once
	err = json.NewDecoder(body).Decode(&sessionResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	//fmt.Printf("%+v\n", sessionResponse)
	return sessionResponse.TokenInfo.Token
}

func SendFax(apiURL string, token string, payload []byte) error {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request (example assumes a POST request)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
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
		return fmt.Errorf("error: %s", resp.Status);
	}

	// Handle the response as needed
	fmt.Println("Request successful!")
	return nil
}

func main() {
	pdfFilePath := "LRT_Map.pdf"
	// Read the entire file into a byte slice
	pdfBytes, err := ioutil.ReadFile(pdfFilePath)
	if err != nil {
		log.Fatal(err)
	}
	// Encode the byte slice to Base64
	encodedString := base64.StdEncoding.EncodeToString(pdfBytes)

	//Check the Base64 encoding
	fmt.Println(encodedString) 

	faxData := FaxData{
	  	Subject: "Test from WEI",
	  	CoverPage: "NONE",
	  	Memo: "Test from Golang App",
	  	DeliveryTime: time.Now(),
	  	Priority: "None",
	}

	//Dynamically adding Recipient(s) and Attachment(s)
	attachment1 := Attachment{Name: "Sample.pdf", Content: encodedString}
	recipient1 := Recipient{FaxNumber: "9783138268", DeliveryType: "Fax"}
	
	faxData.Attachments = append(faxData.Attachments, attachment1)
	faxData.Recipients = append(faxData.Recipients, recipient1)	

	apiURL := "https://ws.biscomfax.com/Session"
	apiURL2 := "https://ws.biscomfax.com/Fax/small"
	username := "cnr_apiuser2"
	password := "Gendarmerie@"

	// Get Session Token
	token := Authenticate(apiURL, username, password)

	if token == "" {
		fmt.Println("Authentication failed")
		return
	}

	payload, err := json.Marshal(faxData)
	if err != nil {
		fmt.Println("Error: error occurred marshalling!")
		return
	}

	// Call SendFax
	err2 := SendFax(apiURL2, token, payload)

	if err2 != nil {
		fmt.Println("Error: error occurred sending fax!")
		return
	}	
}
