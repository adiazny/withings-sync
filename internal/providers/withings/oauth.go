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
	authURL  = "https://account.withings.com/oauth2_user/authorize2"
	tokenURL = "https://wbsapi.withings.net/v2/oauth2"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	CallbackURL  string
	Scopes       []string
}

type provider struct {
	clientID     string
	clientSecret string
	providerName string
	config       Config
	httpClient   *http.Client
	callbackURL  string
}

type refreshBody struct {
	Userid       int    `json:"userid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type refreshResponse struct {
	Status int         `json:"status"`
	Body   refreshBody `json:"body"`
	Error  string      `json:"error"`
}

// NewProvider creates new withings provider
func NewProvider(c *Config) *provider {
	p := new(provider)
	p.providerName = "withings"
	p.clientID = c.ClientID
	p.clientSecret = c.ClientSecret
	p.callbackURL = c.CallbackURL
	p.httpClient = http.DefaultClient

	return p
}

// TODO: Authorize
// TODO: Get Access Token

// RefreshToken refreshes access token
func RefreshToken(p *provider, rt string) (access, refresh string, err error) {
	fmt.Println("Inside RefreshToken() Client ID:", p.clientID)
	formData := url.Values{
		"action":        {"requesttoken"},
		"grant_type":    {"authorization_code"},
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
		"refresh_token": {rt},
	}

	encodedFormData := formData.Encode()
	fmt.Println("Inside RefreshToken() EncodedFormData:", encodedFormData)

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
	fmt.Println(resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var refreshResp refreshResponse
	json.Unmarshal(b, &refreshResp)

	fmt.Printf("RefreshResponse Status: %v\n", refreshResp.Status)
	fmt.Printf("RefreshResponse: %+v\n", refreshResp)
	access = refreshResp.Body.AccessToken
	refresh = refreshResp.Body.RefreshToken

	return

}

// func main() {
// 	fmt.Println("Hello")

// 	c := new(Config)
// 	c.ClientID = os.Getenv("WITHINGS_ID")
// 	c.clientSecret = os.Getenv("WITHINGS_SECRET")

// 	fmt.Printf("client_id %v\n", os.Getenv("WITHINGS_ID"))
// 	fmt.Printf("client_secret %v\n", os.Getenv("WITHINGS_SECRET"))

// 	currentRefreshToken := os.Getenv("WITHINGS_REFRESH")
// 	fmt.Printf("Current Refresh Token: %v\n", currentRefreshToken)

// 	withingsProvider := NewProvider(c)

// 	newAccess, newRefresh, err := RefreshToken(withingsProvider, currentRefreshToken)
// 	if err != nil {
// 		log.Fatalf("Error Refreshing Access Token: %v", err)
// 	}
// 	fmt.Printf("New Access Token: %v\n", newAccess)
// 	fmt.Printf("New Refresh Token: %v\n", newRefresh)

// 	os.Setenv("WITHINGS_REFRESH", newRefresh)
// 	fmt.Printf("New Refresh Env Var: %v\n", os.Getenv("WITHINGS_REFRESH"))
// 	os.Setenv("WITHINGS_TEST", "TEST12345")
// 	fmt.Printf("New TEST Env Var: %v\n", os.Getenv("WITHINGS_TEST"))

// }
