package main

import (
	"bufio"
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

func getMeas(n WithingsNotification) (weight float64) {

	//TODO: Continue testing HTTP client with GET Request

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://account.withings.com/oauth2_user/authorize2", nil)
	if err != nil {
		panic(err)
	}
	qParams := req.URL.Query()
	qParams.Add("response_type", "code")
	qParams.Add("client_id", "3ed034dcb34ebb1c9c2bb300242fca8cdcc3e92a63a1b73b90b7b19413937e60")
	qParams.Add("state", "mystate")
	qParams.Add("scope", "metrics")
	qParams.Add("redirect_uri", "https://portfolio.alandiaz.com")
	req.URL.RawQuery = qParams.Encode()

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
