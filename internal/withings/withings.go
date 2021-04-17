package withings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Notification struct defines body of callback request
type Notification struct {
	UserID    int   `json:"userid"`
	Appli     int   `json:"appli"`
	StartDate int64 `json:"startdate"`
	EndDate   int64 `json:"enddate"`
}

type measures struct {
	Value int `json:"value"`
	Type  int `json:"type"`
	Unit  int `json:"unit"`
}

type measureGrps struct {
	Grpid    int        `json:"grpid"`
	Created  int64      `json:"created"`
	Category int        `json:"category"`
	Measures []measures `json:"measures"`
}

type measureBody struct {
	UpdateTime  int64         `json:"updatetime"`
	TimeZone    string        `json:"timezone"`
	MeasureGrps []measureGrps `json:"measureGrps"`
}

type measureResponse struct {
	Status int         `json:"status"`
	MB     measureBody `json:"body"`
	Error  string      `json:"error"`
}

func NotificationStruct(uv url.Values) (notification Notification) {
	for k, v := range uv {
		switch k {
		case "userid":
			notification.UserID, _ = strconv.Atoi(v[0])
		case "appli":
			notification.Appli, _ = strconv.Atoi(v[0])
		case "startdate":
			notification.StartDate, _ = strconv.ParseInt(v[0], 10, 64)
		case "enddate":
			notification.EndDate, _ = strconv.ParseInt(v[0], 10, 64)
		}
	}
	return
}

// GetMeas calls withings api HTTP POST to get user's weight (kg)
func GetMeas(n Notification, c *Config) (weight float64, err error) {

	wp := NewProvider("withings", c)

	newAccess, _, err := wp.RefreshToken(c.WithingsRefreshToken)
	if err != nil {
		log.Printf("Error Refreshing Access Token: %v", err)
	}

	client := &http.Client{}

	formData := url.Values{
		"action":     {"getmeas"},
		"meastype":   {"1"},
		"lastupdate": {"n.StartDate"},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest("POST", "https://wbsapi.withings.net/measure", strings.NewReader(encodedFormData))
	if err != nil {
		log.Printf("NewRequest Log Err: %v\n", err)
		return weight, fmt.Errorf("Error: %v", err)
	}

	//TODO: Design how to call and add OAuth token
	req.Header.Add("Authorization", "Bearer "+newAccess)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do Log Err: %v\n", err)
		return weight, fmt.Errorf("Error: %v", err)

	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var measResponse measureResponse
	json.Unmarshal(b, &measResponse)

	fmt.Println(measResponse.Status)

	if measResponse.Status != 0 {
		return weight, fmt.Errorf("Error: %v. Status "+strconv.Itoa(measResponse.Status), measResponse.Error)
	}

	weight, err = getWeight(measResponse)
	return
}

func getWeight(mr measureResponse) (weight float64, err error) {
	measuresGrpsList := mr.MB.MeasureGrps
	if len(measuresGrpsList) == 0 {
		return weight, errors.New("No Measures Found")
	}
	lastUpdateMeasures := measuresGrpsList[0].Measures[0]
	weight = calculateWeight(lastUpdateMeasures)
	return
}

func calculateWeight(m measures) (roundedWeight float64) {
	weight := float64(m.Value) * math.Pow10(m.Unit)
	roundedWeight = math.Ceil(weight*100) / 100
	return
}
