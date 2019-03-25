package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/http"
	"github.com/farzadrastegar/simple-cab/driver_location/mock"
	nethttp "net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type patchResponse struct {
	Err string `json:"err,omitempty"`
}

/* Parameters */
// Constants.
const locFormat = "{\"latitude\":%f,\"longitude\":%f,\"updated_at\":\"%s\"}"

var timeNow = time.Date(2019, 03, 16, 00, 00, 00, 00, time.UTC)
var loc1 = driver_location.Location{Latitude: 1.1, Longitude: 1.2, Updated_at: "t1"}
var loc2 = driver_location.Location{Latitude: 2.1, Longitude: 2.2, Updated_at: "t2"}
var loc3 = driver_location.Location{Latitude: 3.1, Longitude: 3.2, Updated_at: "t3"}

// Consumer parameters.
var cLat float64
var cLon float64
var cUpd string
var cId string

/* End of parameters. */

// Mocked outputs.
var storedLocation string
var consumedLocation string

// Services.
var busService driver_location.BusService

func beforeEachTest() {
	driver_location.LoadConfigurationFromBranch()
}

func Test_DriverLocation_StoreLocationWithBus_main(t *testing.T) {
	beforeEachTest()

	s := OpenServer()
	defer s.Close()

	// Prepare sample input.
	cLat = 0.1
	cLon = 0.2
	cUpd = timeNow.Format(time.RFC3339)

	// Consume sample input.
	if err := busService.Consume(); err != nil {
		t.Fatal(err)
	}

	// Make sure consumed message is in right format.
	expectedOut := fmt.Sprintf(locFormat, cLat, cLon, cUpd)
	if expectedOut != consumedLocation {
		t.Fatalf("out=%s, expectedOut=%s", consumedLocation, expectedOut)
	}
}

func Test_DriverLocation_StoreLocationWithHTTP_main(t *testing.T) {
	beforeEachTest()

	s := OpenServer()
	defer s.Close()

	// Prepare sample input.
	cLat = 0.1
	cLon = 0.2
	cUpd = timeNow.Format(time.RFC3339)
	cLoc := &driver_location.Location{Latitude: cLat, Longitude: cLon}
	id := "123"

	// Prepare a PATCH request.
	var u url.URL
	u.Scheme = "HTTP"
	u.Host = fmt.Sprintf(":%d", s.Port())
	u.Path = "/drivers/" + url.QueryEscape(string(id)) + "/locations"
	reqBody, err := json.Marshal(*cLoc)
	if err != nil {
		t.Fatal(err)
	}

	// Create request.
	req, err := nethttp.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Execute request.
	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody patchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatal(err)
	} else if respBody.Err != "" {
		t.Fatal(driver_location.Error(respBody.Err))
	}

	// Compare outputs.
	expectedOut := fmt.Sprintf(locFormat, cLat, cLon, cUpd)
	if expectedOut != storedLocation {
		t.Fatalf("out=%s, expectedOut=%s", storedLocation, expectedOut)
	}
}

func Test_DriverLocation_GetDriverLocationsWithHTTP_main(t *testing.T) {
	beforeEachTest()

	s := OpenServer()
	defer s.Close()

	// Prepare sample input.
	minutes := 1.0
	id := "123"
	locs := &driver_location.Locations{}
	locs.Locations = append(append(append(locs.Locations, loc1), loc2), loc3)
	locs.Err = ""

	// Prepare a GET request.
	var u url.URL
	u.Scheme = "HTTP"
	u.Host = fmt.Sprintf(":%d", s.Port())
	u.Path = "/drivers/" + url.QueryEscape(id) + "/locations"
	u.RawQuery = "minutes=" + fmt.Sprintf("%f", minutes)

	// Create request.
	req, err := nethttp.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute request.
	resp, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody driver_location.Locations
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatal(err)
	}
	respBody.ServerIP = ""

	// Compare outputs.
	if !reflect.DeepEqual(*locs, respBody) {
		t.Fatalf("out=%#v, expectedOut=%#v", respBody, locs)
	}
}

func OpenServer() *http.Server {
	// Mock CabService.
	cabService := &mock.CabService{
		// StoreLocation
		StoreLocationFn: func(id string, loc *driver_location.Location) error {
			storedLocation = fmt.Sprintf(locFormat,
				loc.Latitude,
				loc.Longitude,
				timeNow.Format(time.RFC3339))
			return nil
		},
		//GetDriverLocations
		GetDriverLocationsFn: func(id string, minutes float64) (*driver_location.Locations, error) {
			// Mock 'locations' as the results coming from database.
			loc1 := fmt.Sprintf(locFormat, loc1.Latitude, loc1.Longitude, loc1.Updated_at)
			loc2 := fmt.Sprintf(locFormat, loc2.Latitude, loc2.Longitude, loc2.Updated_at)
			loc3 := fmt.Sprintf(locFormat, loc3.Latitude, loc3.Longitude, loc3.Updated_at)
			locations := fmt.Sprintf("{\"locations\":[%s,%s,%s]}", loc1, loc2, loc3)

			// Decode response into JSON.
			var locStruct driver_location.Locations
			if err := json.NewDecoder(bytes.NewReader([]byte(locations))).Decode(&locStruct); err != nil {
				return nil, err
			}

			return &locStruct, nil
		},
	}

	// Mock BusService.
	busService = &mock.BusService{
		ConsumeFn: func() error {
			consumedLocation = fmt.Sprintf(locFormat, cLat, cLon, cUpd)
			return nil
		},
	}

	// Attach cabService to HTTP handler.
	h := http.NewDataHandler()
	h.CabService = cabService

	// Start an HTTP server.
	srv := http.NewServer()
	srv.Handler = &http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return srv
}
