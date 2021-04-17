package withings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authURL          = "https://account.withings.com/oauth2_user/authorize2"
	withingsTokenURL = "https://wbsapi.withings.net/v2/oauth2"
	stravaTokenURL   = "https://www.strava.com/oauth/token"
)

// Config Struct
type Config struct {
	WithingsClientID     string
	WithingsClientSecret string
	RedirectURL          string
	CallbackURL          string
	Scopes               []string
	WithingsRefreshToken string
	StravaClientID       string
	StravaClientSecret   string
	StravaRefreshToken   string
}

type provider struct {
	providerName string
	config       *Config
	httpClient   *http.Client
}

type refreshBody struct {
	Userid       int    `json:"userid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type withingsRefreshResponse struct {
	Status int         `json:"status"`
	Body   refreshBody `json:"body"`
	Error  string      `json:"error"`
}

type stravaRefreshResponse struct {
	TokenType    string `json:"token_type"`
	AccesToken   string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// NewProvider creates new withings provider
func NewProvider(name string, c *Config) *provider {
	p := new(provider)
	p.providerName = name
	p.config = c
	p.httpClient = http.DefaultClient
	return p
}

// RefreshToken refreshes access token
func (p *provider) RefreshToken(rt string) (access, refresh string, err error) {
	log.Printf("Provider: %s Inside RefreshToken()", p.providerName)

	var tokenURL string
	var formData url.Values

	switch p.providerName {
	case "withings":
		tokenURL = withingsTokenURL
		formData = url.Values{
			"action":        {"requesttoken"},
			"grant_type":    {"authorization_code"},
			"client_id":     {p.config.WithingsClientID},
			"client_secret": {p.config.WithingsClientSecret},
			"refresh_token": {rt},
		}
	case "strava":
		tokenURL = stravaTokenURL
		formData = url.Values{
			"grant_type":    {"refresh_token"},
			"client_id":     {p.config.StravaClientID},
			"client_secret": {p.config.StravaClientSecret},
			"refresh_token": {rt},
		}
	}

	encodedFormData := formData.Encode()

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(encodedFormData))
	if err != nil {
		log.Printf("NewRequest Log Err: %v\n", err)
		return access, refresh, fmt.Errorf("Error: %v", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Printf("client.Do Log Err: %v\n", err)
		return access, refresh, fmt.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	switch p.providerName {
	case "withings":
		var withingsRefreshResp withingsRefreshResponse
		json.Unmarshal(b, &withingsRefreshResp)
		log.Printf("RefreshResponse: %+v\n", withingsRefreshResp)
		access = withingsRefreshResp.Body.AccessToken
		refresh = withingsRefreshResp.Body.RefreshToken

		return
	case "strava":
		var stravaRefreshResp stravaRefreshResponse
		json.Unmarshal(b, &stravaRefreshResp)
		log.Printf("RefreshResponse: %+v\n", stravaRefreshResp)
		access = stravaRefreshResp.AccesToken
		refresh = stravaRefreshResp.RefreshToken
		return
	}
	return

}
