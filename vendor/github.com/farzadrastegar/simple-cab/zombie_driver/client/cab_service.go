package client

import (
	"encoding/json"
	"fmt"
	"github.com/alouche/go-geolib"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"net/http"
	"net/url"
	"strconv"
)

//const clientYamlFilename = "../config.yaml"

var localServerAddr = "localhost"
var localServerPort = "8082"
var distance = "500.0"
var distanceFloat = 500.0
var duration = "5.0"

func GetDuration() string {
	return duration
}

// Ensure service implements interface.
var _ zombie_driver.CabService = &CabService{}

type CabService struct {
	//client *Client
	handler **Handler
}

func (s *CabService) CheckZombieStatus(id string) (*zombie_driver.Status, error) {
	var u url.URL
	u.Path = "/drivers/" + url.QueryEscape(id) + "/locations"
	u.RawQuery = "minutes=" + duration
	u.Scheme = "HTTP"
	u.Host = fmt.Sprintf("%s:%s", localServerAddr, localServerPort)

	// Prepare request.
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Execute request.
	h := *s.handler
	resp, err := h.RequestService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody GetDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Err != "" {
		return nil, zombie_driver.Error(respBody.Err)
	}

	// Calculate distance.
	drivenDistance := calculateDistance(respBody.Locations)

	// Set output elements.
	zombieStatus := false
	if drivenDistance < distanceFloat {
		zombieStatus = true
	}

	// Send output.
	idInt, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		return nil, err
	}
	return &zombie_driver.Status{ID: zombie_driver.DriverID(idInt), Zombie: zombieStatus}, nil
}

func calculateDistance(locations []SingleGeoRecord) float64 {
	locationsLen := len(locations)
	distance := 0.0

	if locationsLen < 2 {
		return distance
	}

	for i := 0; i < locationsLen-1; i++ {
		distance += geolib.HaversineDistance(
			locations[i].Latitude,
			locations[i].Longitude,
			locations[i+1].Latitude,
			locations[i+1].Longitude)
	}

	return distance * 1000.0 //return calculated distance in meters
}

type SingleGeoRecord struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Updated_at string  `json:"updated_at"`
}

type GetDataResponse struct {
	Locations []SingleGeoRecord `json:"locations,omitempty"`
	Err       string            `json:"err,omitempty"`
}

var _ RequestService = &ReqService{}

type ReqService struct{}

func (s *ReqService) ExecuteRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
