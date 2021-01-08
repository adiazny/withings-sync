package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WithingsNotification struct {
	UserID    int   `json:"userid"`
	Appli     int   `json:"appli"`
	StartDate int64 `json:"startdate"`
	EndDate   int64 `json:"enddate"`
}

type Measures struct {
	Value int `json:"value"`
	Type  int `json:"type"`
	Unit  int `json:"unit"`
}

type MeasureGrps struct {
	Grpid    int        `json:"grpid"`
	Created  int64      `json:"created"`
	Category int        `json:"category"`
	Measures []Measures `json:"measures"`
}

type MeasureBody struct {
	UpdateTime  int64         `json:"updatetime"`
	TimeZone    string        `json:"timezone"`
	MeasureGrps []MeasureGrps `json:"measuregrps"`
}

type MeasureResponse struct {
	Status int         `json:"status"`
	MB     MeasureBody `json:"body"`
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
	var notification WithingsNotification
	err = json.Unmarshal(b, &notification)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func convertToTime(unixTime int64) time.Time {
	time := time.Unix(unixTime, 0)
	return time
}

func getMeas(n WithingsNotification) (weight float64) {

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
	req.Header.Add("Authorization", "Bearer 0ff90340b36abcc09df998d1e287669b5f5c8b34")

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
	var meas MeasureResponse
	json.Unmarshal(b, &meas)

	fmt.Printf("Measurement %+v", meas)

	weight = 1.0
	return
}

func main() {
	fmt.Println("Server Up")

	http.HandleFunc("/about", aboutHandler)

	http.HandleFunc("/weight", withingsNotificationHandler)

	fmt.Println(getMeas(WithingsNotification{}))

	log.Fatal(http.ListenAndServe(":8090", nil))

}
