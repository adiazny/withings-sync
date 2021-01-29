package api

import (
	"fmt"
	"net/http"

	"github.com/adiazny/withings-sync/internal/strava"
	"github.com/adiazny/withings-sync/internal/withings"
)

// AboutHandler is a basic endpoint to test liveness
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About page using ResponseWriter.Write for the Withings Sync Application\n"))
	fmt.Fprintf(w, "About page using fmt.Fprintf for the Withings Sync Application.\nRequest Path: %v", r.URL.Path[1:])
}

// WithingsNotificationHandler will handle the HTTP Requests
func WithingsNotificationHandler(w http.ResponseWriter, r *http.Request) {

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
		weight, err := withings.GetMeas(notification)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Update weight in Strava
		err = strava.UpdateWeight(weight)
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

func handleNotificationRequest(r *http.Request) (notification withings.Notification, jsonError error) {

	//Parse Form params
	r.ParseForm()
	urlValues := r.Form

	// Create WithingsNotification Struct from Form params
	notification = withings.NotificationStruct(urlValues)
	return

}
