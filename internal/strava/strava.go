package strava

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/adiazny/withings-sync/internal/withings"
)

// UpdateWeight will perform a PUT request to update user's weight
func UpdateWeight(weight float64, c *withings.Config) error {
	sp := withings.NewProvider("strava", c)
	var err error
	urlString := "https://www.strava.com/api/v3/athlete?weight="
	weightString := strconv.FormatFloat(weight, 'f', 2, 64)

	client := &http.Client{}

	req, err := http.NewRequest("PUT", urlString+weightString, strings.NewReader(""))
	if err != nil {
		log.Printf("NewRequest Log Err: %v\n", err)
		return fmt.Errorf("Error: %v", err)
	}

	access, _, _ := sp.RefreshToken(c.StravaRefreshToken)
	bearer := fmt.Sprintf("Bearer %s", access)
	req.Header.Add("Authorization", bearer)

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
