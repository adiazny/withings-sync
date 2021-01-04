package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type WithingsNotification struct {
	UserID    int   `json:"userid"`
	Appli     int   `json:"appli"`
	StartDate int64 `json:"startdate"`
	EndDate   int64 `json:"enddate"`
}

func main() {
	fmt.Println("Server Up")
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/weight", withingsNotificationHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))

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
		log.Fatal(err)
	}

	//Print Notification
	fmt.Printf("Notification: %+v\n", notification)

	//Convert to Time and Print
	startTime := convertToTime(notification.StartDate)
	endTime := convertToTime(notification.EndDate)
	fmt.Printf("Start Time: %v\n", startTime)
	fmt.Printf("End Time: %v\n", endTime)

}

func convertToTime(unixTime int64) time.Time {
	time := time.Unix(unixTime, 0)
	return time
}
