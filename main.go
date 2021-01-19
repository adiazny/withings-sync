package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type withingsNotification struct {
	UserID    int   `json:"userid"`
	Appli     int   `json:"appli"`
	StartDate int64 `json:"startdate"`
	EndDate   int64 `json:"enddate"`
}

type measures struct {
	Value int `json:"value"`
	Type  int `json:"type"`
	Unit  int `json:"unit"`
}

type measureGrps struct {
	Grpid    int        `json:"grpid"`
	Created  int64      `json:"created"`
	Category int        `json:"category"`
	Measures []measures `json:"measures"`
}

type measureBody struct {
	UpdateTime  int64         `json:"updatetime"`
	TimeZone    string        `json:"timezone"`
	MeasureGrps []measureGrps `json:"measureGrps"`
}

type measureResponse struct {
	Status int         `json:"status"`
	MB     measureBody `json:"body"`
	Error  string      `json:"error"`
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About page using ResponseWriter.Write for the Withings Sync Application\n"))
	fmt.Fprintf(w, "About page using fmt.Fprintf for the Withings Sync Application.\nRequest Path: %v", r.URL.Path[1:])
}

func withingsNotificationHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		fmt.Fprintln(w, "Get Request Successful")
	case "POST":
		notification, err := handleNotificationRequest(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Pull weight from Withings API
		weight, err := getMeas(notification)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Update weight in Strava
		err = updateStravaWeight(weight)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "Weight: %v", weight)
		fmt.Println("User Weight (kg):", weight)
	default:
		fmt.Fprintf(w, "Only GET, HEAD and POST allowed")
	}

}

func handleNotificationRequest(r *http.Request) (notification withingsNotification, jsonError error) {

	//Parse Form params
	r.ParseForm()
	urlValues := r.Form

	// Create WithingsNotification Struct from Form params
	notification = notificationStruct(urlValues)
	return

}

func notificationStruct(uv url.Values) (notification withingsNotification) {
	for k, v := range uv {
		switch k {
		case "userid":
			notification.UserID, _ = strconv.Atoi(v[0])
		case "appli":
			notification.Appli, _ = strconv.Atoi(v[0])
		case "startdate":
			notification.StartDate, _ = strconv.ParseInt(v[0], 10, 64)
		case "enddate":
			notification.EndDate, _ = strconv.ParseInt(v[0], 10, 64)
		}
	}
	return
}

func convertToTime(unixTime int64) time.Time {
	time := time.Unix(unixTime, 0)
	return time
}

func getMeas(n withingsNotification) (weight float64, err error) {
	client := &http.Client{}

	formData := url.Values{
		"action":     {"getmeas"},
		"meastype":   {"1"},
		"lastupdate": {"n.StartDate"},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest("POST", "https://wbsapi.withings.net/measure", strings.NewReader(encodedFormData))
	if err != nil {
		log.Printf("NewRequest Log Err: %v\n", err)
		return weight, fmt.Errorf("Error: %v", err)
	}

	//TODO: Design how to call and add OAuth token
	req.Header.Add("Authorization", "Bearer XXXX")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do Log Err: %v\n", err)
		return weight, fmt.Errorf("Error: %v", err)

	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var measResponse measureResponse
	json.Unmarshal(b, &measResponse)

	fmt.Println(measResponse.Status)

	if measResponse.Status != 0 {
		return weight, fmt.Errorf("Error: %v. Status "+strconv.Itoa(measResponse.Status), measResponse.Error)
	}

	weight, err = getWeight(measResponse)
	return
}

func getWeight(mr measureResponse) (weight float64, err error) {
	measuresGrpsList := mr.MB.MeasureGrps
	if len(measuresGrpsList) == 0 {
		return weight, errors.New("No Measures Found")
	}
	lastUpdateMeasures := measuresGrpsList[0].Measures[0]
	weight = calculateWeight(lastUpdateMeasures)
	return
}

func calculateWeight(m measures) (roundedWeight float64) {
	weight := float64(m.Value) * math.Pow10(m.Unit)
	roundedWeight = math.Ceil(weight*100) / 100
	return
}

func updateStravaWeight(weight float64) error {
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

func main() {
	fmt.Println("Server Up")

	http.HandleFunc("/about", aboutHandler)

	http.HandleFunc("/weight", withingsNotificationHandler)

	log.Fatal(http.ListenAndServe(":8090", nil))

}
