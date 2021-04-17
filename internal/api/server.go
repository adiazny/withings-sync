package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adiazny/withings-sync/internal/strava"
	"github.com/adiazny/withings-sync/internal/withings"
)

// Server Struct
type Server struct {
	router *http.ServeMux
	config *withings.Config
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleAboutEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("About page using ResponseWriter.Write for the Withings Sync Application\n"))
	}
}

func (s *Server) handleWithingsCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprintln(w, "Get Request Successful")
		case "POST":
			notification, err := extractFormParams(r)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			// Pull weight from Withings API
			weight, err := withings.GetMeas(notification, s.config)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			// Update weight in Strava
			err = strava.UpdateWeight(weight, s.config)
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
}

// NewServer returns a new Server value
func NewServer() *Server {
	s := &Server{
		router: http.NewServeMux(),
		config: loadConfig(),
	}
	s.routes()

	log.Printf("WITHINGS_REFRESH_TOKEN= %s", s.config.WithingsRefreshToken)
	log.Printf("STRAVA_REFRESH_TOKEN= %s", s.config.StravaRefreshToken)
	log.Printf("STRAVA_CLIENT_ID= %s", s.config.StravaClientID)
	log.Println("Withings-Sync Application Started Successfully")
	return s
}

func loadConfig() *withings.Config {

	log.Println("Loading configuration from env vars.")
	return &withings.Config{
		WithingsClientID:     getEnvVar("WITHINGS_CLIENT_ID"),
		WithingsClientSecret: getEnvVar("WITHINGS_CLIENT_SECRET"),
		RedirectURL:          getEnvVar("WITHINGS_REDIRECT"),
		CallbackURL:          getEnvVar("WITHINGS_CALLBACK"),
		WithingsRefreshToken: getEnvVar("WITHINGS_REFRESH_TOKEN"),
		StravaClientID:       getEnvVar("STRAVA_CLIENT_ID"),
		StravaClientSecret:   getEnvVar("STRAVA_CLIENT_SECRET"),
		StravaRefreshToken:   getEnvVar("STRAVA_REFRESH_TOKEN"),
	}
}

func getEnvVar(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Environment variable %v is not set", key)
	}
	if val == "" {
		log.Fatalf("Environment variable %v is empty", key)
	}
	return val
}

func extractFormParams(r *http.Request) (notification withings.Notification, jsonError error) {

	//Parse Form params
	r.ParseForm()
	urlValues := r.Form

	// Create WithingsNotification Struct from Form params
	notification = withings.NotificationStruct(urlValues)
	return

}
