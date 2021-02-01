package main

import (
	"log"
	"net/http"

	"github.com/adiazny/withings-sync/internal/api"
)

func main() {
	log.Println("Starting Withings-Sync Application...")

	log.Fatal(http.ListenAndServe(":8090", api.NewServer()))
}
