package strava

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// UpdateWeight will make strava api HTTP PUT to update user's weight
func UpdateWeight(weight float64) error {
	var err error
	urlString := "https://www.strava.com/api/v3/athlete?weight="
	weightString := strconv.FormatFloat(weight, 'f', 2, 64)

	client := &http.Client{}

	req, err := http.NewRequest("PUT", urlString+weightString, strings.NewReader(""))
	if err != nil {
		log.Printf("NewRequest Log Err: %v\n", err)
		return fmt.Errorf("Error: %v", err)
	}
	//TODO: Design how to call and add OAuth token
	req.Header.Add("Authorization", "Bearer XXXX")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do Log Err: %v\n", err)
		return fmt.Errorf("Error: %v", err)

	}
	defer resp.Body.Close()

	log.Printf("Strava PUT Response Status: %v", resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Strava Error: %v", resp.Status)
	}

	return err

}
