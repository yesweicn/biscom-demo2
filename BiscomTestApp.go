package main

import (
	"BiscomLibrary"	
	"encoding/json"
	"fmt"
)

import 

func main() {
	faxData := BiscomLibrary.FaxData {
	  	subject: "Test from Biscom",
	  	coverPage: "NONE",
	  	memo: "Test from Me",
	  	deliveryTime: time.Now(),
	  	priority: "None",
	  	recipients: []Recipient{
			{
		  	"faxNumber": "9783136067",
		  	"deliveryType": "Fax"
			}
	  	}
	}

	apiURL := "https://ws.biscomfax.com/Session"
	apiURL2 := "https://ws.biscomfax.com/Fax/small"
	username := "your_username"
	password := "your_password"

	// Call OperationA
	token := BiscomLibrary.Autenticate(apiURL, username, password)

	payload := json.Marshal(faxData)

	// Call OperationB
	Biscomlibrary.SendFax(apiURL, token, payload)
}
