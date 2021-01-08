package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
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
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About page using ResponseWriter.Write for the Withings Sync Application\n"))
	fmt.Fprintf(w, "About page using fmt.Fprintf for the Withings Sync Application.\nRequest Path: %v", r.URL.Path[1:])
}

func withingsNotificationHandler(w http.ResponseWriter, r *http.Request) {
	//Read Body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Unmarshal
	var notification withingsNotification
	err = json.Unmarshal(b, &notification)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	weight := getMeas(notification)

	fmt.Fprintf(w, "Weight: %v", weight)

	fmt.Println("User Weight (kg):", weight)

}

func convertToTime(unixTime int64) time.Time {
	time := time.Unix(unixTime, 0)
	return time
}

func getMeas(n withingsNotification) (weight float64) {

	client := &http.Client{}

	formData := url.Values{
		"action":    {"getmeas"},
		"meastype":  {"1"},
		"starttime": {"n.StartDate"},
		"endtime":   {"n.EndDate"},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest("POST", "https://wbsapi.withings.net/measure", strings.NewReader(encodedFormData))
	if err != nil {
		panic(err)
	}
	//TODO: Design how to call and add OAuth token
	req.Header.Add("Authorization", "Bearer a4df0711085010b47f25e6ade4820f54ed8f08e7")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response status:", resp.Status)

	//var jsonData interface{}
	var meas measureResponse
	json.Unmarshal(b, &meas)

	measuresGrpsList := meas.MB.MeasureGrps
	sort.Slice(measuresGrpsList, func(i, j int) bool {
		return measuresGrpsList[i].Created > measuresGrpsList[j].Created
	})

	recentMeasures := measuresGrpsList[0].Measures[0]

	weight = calculateWeight(recentMeasures)
	return
}

func calculateWeight(m measures) (weight float64) {
	// TODO: look up how to round to nearest tenth
	return float64(m.Value) * math.Pow10(m.Unit)
}

func main() {
	fmt.Println("Server Up")

	http.HandleFunc("/about", aboutHandler)

	http.HandleFunc("/weight", withingsNotificationHandler)

	log.Fatal(http.ListenAndServe(":8090", nil))

}
