package http_test

import (
	"bytes"
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/http"
	"github.com/farzadrastegar/simple-cab/gateway/mock"
	"github.com/julienschmidt/httprouter"
	"log"
	nethttp "net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestCabService_StoreLocation_InvalidInput(t *testing.T) {
	t.Run("StatusBadRequest", testCabService_StoreLocation_InvalidInputJson)
	t.Run("MethodNotAllowed", testCabService_StoreLocation_InvalidHTTPMethod)
}

// Ensure service returns right status code for invalid input json.
func testCabService_StoreLocation_InvalidInputJson(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	c.cabService.StoreLocationFn = func(id string, data *gateway.Data) error {
		u := c.URL
		u.Path = "/drivers/" + url.QueryEscape(string(id)) + "/locations"

		// Create a router and setup routes
		h := http.NewDataHandler()

		// Create an invalid request body.
		reqBodyStr := "{latitude:1.234567, longitude:1.234567"
		reqBody := []byte(reqBodyStr)

		// Create request.
		req, err := nethttp.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
		if err != nil {
			return err
		}

		// Execute request.
		w := testHTTPResponse(t, h.Router, req)
		return gateway.Error(w.Result().StatusCode)
	}

	// Send a request
	err := c.Connect().StoreLocation("123", nil)

	if err.Error() != string(nethttp.StatusBadRequest) {
		t.Fatal(err)
	}
}

// Ensure service returns right status code for invalid method.
func testCabService_StoreLocation_InvalidHTTPMethod(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	c.cabService.StoreLocationFn = func(id string, data *gateway.Data) error {
		u := c.URL
		u.Path = "/drivers/" + url.QueryEscape(string(id)) + "/locations"

		// Create a router and setup routes
		h := http.NewDataHandler()

		// Create an invalid request body.
		reqBodyStr := "{latitude:1.234567, longitude:1.234567}"
		reqBody := []byte(reqBodyStr)

		// Create request.
		req, err := nethttp.NewRequest("POST", u.String(), bytes.NewReader(reqBody))
		if err != nil {
			return err
		}

		// Execute request.
		w := testHTTPResponse(t, h.Router, req)
		return gateway.Error(w.Result().StatusCode)
	}

	// Send a request
	err := c.Connect().StoreLocation("123", nil)

	if err.Error() != string(nethttp.StatusMethodNotAllowed) {
		t.Fatal(err)
	}
}

func TestCabService_CheckValidOutput(t *testing.T) {
	t.Run("StoreLocation", testCabService_StoreLocation_ValidOuput)
	t.Run("CheckZombieStatus", testCabService_CheckZombieStatus_ValidOutput)
}

// Ensure service returns right output.
func testCabService_StoreLocation_ValidOuput(t *testing.T) {
	s, c := MustOpenTestServerHttpClient()
	defer s.Close()

	// Mock service.
	s.Handler.DataHandler.CabService.StoreLocationFn = func(id string, data *gateway.Data) error {
		return nil
	}

	// Create sample data.
	d := gateway.Data{Latitude: 1.234567, Longitude: 1.234567}

	// Send a request.
	err := c.Connect().StoreLocation("123", &d)

	if err != nil {
		t.Fatal(err)
	}
}

// Ensure service returns right output.
func testCabService_CheckZombieStatus_ValidOutput(t *testing.T) {
	s, c := MustOpenTestServerHttpClient()
	defer s.Close()

	// Create sample data.
	expectedOut := gateway.Status{ID: 123, Zombie: false}

	// Mock service.
	s.Handler.DataHandler.CabService.CheckZombieStatusFn = func(id string) (*gateway.Status, error) {
		return &expectedOut, nil
	}

	// Send a request.
	out, err := c.Connect().CheckZombieStatus("123")

	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(expectedOut, *out) {
		t.Fatalf("Output:%#v, Expected:%#v", *out, expectedOut)
	}
}

// testHTTPResponse is a helper function to process a request and return its response.
func testHTTPResponse(t *testing.T, r *httprouter.Router, req *nethttp.Request) *httptest.ResponseRecorder {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	return w
}

// DataHandler represents a test wrapper for http.DataHandler.
type DataHandler struct {
	*http.DataHandler

	CabService mock.CabService
	LogOutput  bytes.Buffer
}

// NewDataHandler returns a new instance of DataHandler.
func NewDataHandler() *DataHandler {
	h := &DataHandler{DataHandler: http.NewDataHandler()}
	h.DataHandler.CabService = &h.CabService
	h.Logger = log.New(VerboseWriter(&h.LogOutput), "", log.LstdFlags|log.Lshortfile)
	return h
}
