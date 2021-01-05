package main

import (
	"bufio"
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
	req.Header.Add("Authorization", "Bearer e43d13aedc108b06aa259fe3a4a2acaa70caedfa")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
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
