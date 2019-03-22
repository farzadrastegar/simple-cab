package client_test

import (
	"encoding/json"
	"fmt"
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/client"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

// Ensure service stores data into a bus.
func TestCabService_StoreLocation(t *testing.T) {
	// Create client.
	c := NewClient()

	// Mock BusService
	formatOut := "id:%d,{latitude:%f,longitude:%f}"
	var out string
	c.Handler.BusService.ProduceFn = func(id gateway.DriverID, message *gateway.Data) error {
		out = fmt.Sprintf(formatOut, id, message.Latitude, message.Longitude)
		return nil
	}

	// Prepare sample data.
	id := "123"
	idInt := 123
	data := gateway.Data{Latitude: 1.234567, Longitude: 1.234567}

	// Store data.
	expectedOut := fmt.Sprintf(formatOut, idInt, data.Latitude, data.Longitude)
	err := c.Connect().StoreLocation(id, &data)

	if err != nil {
		t.Fatal(err)
	} else if out != expectedOut {
		t.Fatalf("out=%s, expectedOut=%s", out, expectedOut)
	}
}

// Ensure service returns right output.
func TestCabService_CheckZombieStatus(t *testing.T) {
	// Create client.
	c := NewClient()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	// Prepare sample data.
	id := "123"
	idInt := 123
	expectedStatus := gateway.Status{ID: gateway.DriverID(idInt), Zombie: false}

	// Mock service.
	c.Handler.RequestService.r = httprouter.New()
	c.Handler.RequestService.r.GET("/drivers/:id", func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		//encodeJSON(writer, &client.GetDataResponse{Err: "", Status: &expectedStatus}, logger)
		encodeJSON(writer, &gateway.Status{ID: expectedStatus.ID, Zombie: expectedStatus.Zombie}, logger)
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
		err = gateway.ErrInternal
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

	// Create the service and process the request.
	s.r.ServeHTTP(w, req)

	return w.Result(), nil
}
