package client_test

import (
	"encoding/json"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"github.com/farzadrastegar/simple-cab/zombie_driver/client"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

// Ensure service returns right output.
func TestCabService_CheckZombieStatus(t *testing.T) {
	// Create client.
	c := NewClient()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	// Prepare sample data.
	id := "123"
	idInt := 123
	expectedStatus := zombie_driver.Status{ID: zombie_driver.DriverID(idInt), Zombie: false}

	// Mock service.
	c.Handler.RequestService.r = httprouter.New()
	c.Handler.RequestService.r.GET("/drivers/:id/locations?minutes="+client.GetDuration(), func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		encodeJSON(writer, &expectedStatus, logger)
	})

	// Call service.
	response, err := c.Connect().CheckZombieStatus(id)

	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(expectedStatus, *response) {
		t.Fatalf("out=%#v, expectedOut=%#v", *response, expectedStatus)
	}
}

func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = zombie_driver.ErrInternal
	}

	// Write generic error response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

// errorResponse is a generic response for sending a error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

var _ client.RequestService = &RequestService{}

type RequestService struct {
	r *httprouter.Router
}

func (s *RequestService) ExecuteRequest(req *http.Request) (*http.Response, error) {
	// Create a response recorder.
	w := httptest.NewRecorder()

	//// Create the service and process the request.
	//s.r.ServeHTTP(w, req)

	// Prepare output.
	elem1 := client.SingleGeoRecord{Latitude: 1.2, Longitude: 1.2}
	elem2 := client.SingleGeoRecord{Latitude: 2.2, Longitude: 2.2}
	elem3 := client.SingleGeoRecord{Latitude: 3.2, Longitude: 3.2}
	loc := client.GetDataResponse{Locations: []client.SingleGeoRecord{elem1, elem2, elem3}, Err: ""}
	json.NewEncoder(w).Encode(&loc)
	return w.Result(), nil
}
