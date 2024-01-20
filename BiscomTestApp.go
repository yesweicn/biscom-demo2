package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	faxData := FaxData{
	  	Subject: "Test from Biscom",
	  	CoverPage: "NONE",
	  	Memo: "Test from Me",
	  	DeliveryTime: time.Now(),
	  	Priority: "None",
	  	Recipients: []Recipient{
			{
		  		FaxNumber: "9783136067",
		  		DeliveryType: "Fax",
			},
	  	},
	}

	apiURL := "https://ws.biscomfax.com/Session"
	apiURL2 := "https://ws.biscomfax.com/Fax/small"
	username := "cnr_apiuser"
	password := "Cnr_7806040871"

	// Call OperationA
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
