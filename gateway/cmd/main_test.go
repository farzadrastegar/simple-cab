package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/client"
	"github.com/farzadrastegar/simple-cab/gateway/config"
	"github.com/farzadrastegar/simple-cab/gateway/http"
	"github.com/farzadrastegar/simple-cab/gateway/mock"
	"log"
	nethttp "net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var busOut string

const busOutFormat = "id:%d,{latitude:%f,longitude:%f}"

func beforeEachTest() {
	// Set yaml filename.
	gateway.SetConfigFilename("../cmd/config.yaml")
}

func Test_Gateway_StoreLocation_main(t *testing.T) {
	beforeEachTest()

	s := OpenServer()
	defer s.Close()

	// Create sample input.
	id := "123"
	idInt, _ := strconv.ParseInt(id, 0, 64)
	locInput := gateway.Data{Latitude: 1.234567, Longitude: 1.234567}
	expectedOut := fmt.Sprintf(busOutFormat, idInt, locInput.Latitude, locInput.Longitude)

	// Encode request body.
	reqBody, err := json.Marshal(locInput)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Patch request.
	serverAddrPort := fmt.Sprintf("localhost:%d", s.Port())
	u := url.URL{Scheme: "HTTP", Host: serverAddrPort, Path: "/drivers/" + id + "/locations"}
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

	if expectedOut != busOut {
		t.Fatalf("Out:%s, Expected out=%s", busOut, expectedOut)
	}
}

func Test_Gateway_CheckZombieStatus_main(t *testing.T) {
	beforeEachTest()

	// Read CheckZombieStatus internal service's address and port.
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	conf := config.NewConfig(logger)
	conf.ReadYaml(gateway.GetConfigFilename())
	zPort := conf.GetYamlValueStr("servers", "internal", "port")

	// Check CheckZombieStatus port is listening
	checkPort1 := fmt.Sprintf("lsof -i -n -P | grep %s | grep LISTEN | tail -n1", zPort)
	cmdOut, _ := exec.Command("/bin/bash", "-c", checkPort1).Output()
	var strBuilder strings.Builder
	strBuilder.Write(cmdOut)
	if strings.Index(strBuilder.String(), zPort) == -1 {
		t.Skipf("internal port %s is not open", zPort)
	}

	// Open server.
	s := OpenServer()
	defer s.Close()

	// Create sample input/output data.
	id := "123"
	idInt, _ := strconv.ParseInt(id, 0, 64)
	expectedOut := gateway.Status{ID: gateway.DriverID(idInt), Zombie: false}

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
	var out gateway.Status
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, expectedOut) {
		t.Fatalf("Out:%#v, Expected out=%#v", out, expectedOut)
	}
}

func OpenServer() *http.Server {
	// Mock BusService.
	b := mock.BusService{}
	b.ProduceFn = func(id gateway.DriverID, message *gateway.Data) error {
		busOut = fmt.Sprintf(busOutFormat, id, message.Latitude, message.Longitude)
		return nil
	}

	// Create a client and attach bus service to it.
	c := client.NewClient()
	c.Handler = client.NewHandler()
	c.Handler.BusService = &b

	// Create cab services.
	s := gateway.CreateCabService(c)

	// Attach services to HTTP handler.
	h := http.NewDataHandler()
	h.CabService = s

	// Start an HTTP server.
	srv := http.NewServer()
	srv.Handler = &http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return srv
}
