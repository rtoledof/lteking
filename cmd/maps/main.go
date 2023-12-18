package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	
	// Create the URL for the directions API.
	url := "https://api.mapbox.com/directions/v5/mapbox/driving/uuid/sYRsjBf8PEJiaOoklAttghQ21o6MaTxbXSUvhQkG8_sztca7DgopYw==?access_token=pk.eyJ1IjoicnRvbGVkb2Zlcm5hbmRleiIsImEiOiJjbHB1aTM1ZXMwbTYxMmlxc3FpYTRmdXd0In0.hIytBvseVldGc-tioC-7_Q"

	// Create the HTTP request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check the response status code.
	if resp.StatusCode != 200 {
		fmt.Println("Error:", resp.Status)
		return
	}

	// Decode the response body.
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the response body.
	fmt.Println(response)
}
