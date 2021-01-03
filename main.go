package main

import (
	"fmt"
	"log"
	"net/http"
)

type WithingsMessage struct {
}

func main() {

	http.HandleFunc("/about", aboutHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))

}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About page using ResponseWriter.Write for the Withings Sync Application\n"))
	fmt.Fprintf(w, "About page using fmt.Fprintf for the Withings Sync Application.\nRequest Path: %v", r.URL.Path[1:])
}

func withingsNotificationHandler(w http.ResponseWriter, r *http.Request) {
}
