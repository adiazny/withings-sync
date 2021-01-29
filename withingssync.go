package main

import (
	"log"
	"net/http"

	"github.com/adiazny/withings-sync/internal/api"
)

func main() {
	log.Println("Starting Withings-Sync Application...")

	http.HandleFunc("/about", api.AboutHandler)

	http.HandleFunc("/callback", api.WithingsNotificationHandler)

	log.Fatal(http.ListenAndServe(":8090", nil))

}
