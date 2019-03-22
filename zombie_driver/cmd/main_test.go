package main_test

import (
	"encoding/json"
	"fmt"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"github.com/farzadrastegar/simple-cab/zombie_driver/http"
	"github.com/farzadrastegar/simple-cab/zombie_driver/mock"
	nethttp "net/http"
	"net/url"
	"reflect"
	"testing"
)

// Create sample data.
var idInt = 123
var expectedOut = zombie_driver.Status{ID: zombie_driver.DriverID(idInt), Zombie: false}

func beforeEachTest() {
	// Set yaml filename.
	zombie_driver.SetConfigFilename("../cmd/config.yaml")
}

func Test_ZombieDriver_main(t *testing.T) {
	beforeEachTest()

	// Open server.
	s := OpenServer()
	defer s.Close()

	// Create sample input/output data.
	id := string(idInt)

	// Create a Get request.
	serverAddrPort := fmt.Sprintf("localhost:%d", s.Port())
	u := url.URL{Scheme: "HTTP", Host: serverAddrPort, Path: "/drivers/" + id}
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
	var out zombie_driver.Status
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, expectedOut) {
		t.Fatalf("Out:%#v, Expected out=%#v", out, expectedOut)
	}
}

func OpenServer() *http.Server {
	// Mock CabService.
	s := mock.CabService{}
	s.CheckZombieStatusFn = func(id string) (*zombie_driver.Status, error) {
		return &expectedOut, nil
	}

	// Attach services to HTTP handler.
	h := http.NewDataHandler()
	h.CabService = &s

	// Start an HTTP server.
	srv := http.NewServer()
	srv.Handler = &http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return srv
}
