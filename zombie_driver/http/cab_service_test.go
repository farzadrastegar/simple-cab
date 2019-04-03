package http_test

import (
	"bytes"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"github.com/farzadrastegar/simple-cab/zombie_driver/http"
	"github.com/farzadrastegar/simple-cab/zombie_driver/mock"

	"github.com/julienschmidt/httprouter"
)

func TestCabService_CheckValidOutput(t *testing.T) {
	zombie_driver.LoadConfigurationFromBranch()

	t.Run("CheckZombieStatus", testCabService_CheckZombieStatus_ValidOutput)
}

// Ensure service returns right output.
func testCabService_CheckZombieStatus_ValidOutput(t *testing.T) {
	s, c := MustOpenTestServerHttpClient()
	defer s.Close()

	// Create sample data.
	id := 123
	expectedOut := zombie_driver.Status{ID: zombie_driver.DriverID(id), Zombie: false}

	// Mock service.
	s.Handler.DataHandler.CabService.CheckZombieStatusFn = func(id string) (*zombie_driver.Status, error) {
		return &expectedOut, nil
	}

	// Send a request.
	out, err := c.Connect().CheckZombieStatus(string(id))

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
	return h
}
