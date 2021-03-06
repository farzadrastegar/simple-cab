package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/farzadrastegar/simple-cab/driver_location"

	"github.com/julienschmidt/httprouter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ErrInvalidJSON = driver_location.Error("invalid json")

//const driverLocationConfigFilename = "../config.yaml"
const queryName = "minutes"

// DataHandler represents an HTTP API handler for requests.
type DataHandler struct {
	*httprouter.Router

	CabService driver_location.CabService
}

// NewDataHandler returns a new instance of DataHandler.
func NewDataHandler() *DataHandler {
	h := &DataHandler{
		Router: httprouter.New(),
	}

	//setup routes
	h.SetupRoutes()

	return h
}

func (h *DataHandler) SetupRoutes() {
	//check method types in yaml
	methodType := viper.GetString("urls.driverLocations.method")
	if !strings.EqualFold(methodType, http.MethodGet) &&
		!strings.EqualFold(methodType, http.MethodPatch) {
		logger.Fatalf("ERROR: wrong method type in yaml config file")
	}

	//set routes of endpoints
	h.PATCH(viper.GetString("urls.driverLocations.path"), h.StoreLocation)
	h.GET(viper.GetString("urls.driverLocations.path"), h.GetDriverLocations)
}

//Optional
func (h *DataHandler) StoreLocation(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	// Decode request.
	var req driver_location.Location
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		Error(writer, ErrInvalidJSON, http.StatusBadRequest)
		return
	} else {
		//get URL parameters
		driverId := ps.ByName("id")

		//send request body via NSQ
		switch err := h.CabService.StoreLocation(driverId, &req); err {
		case nil:
			encodeJSON(writer, &patchResponse{})
		default:
			Error(writer, driver_location.ErrInternal, http.StatusBadRequest)
		}
	}
}

type patchResponse struct {
	Err string `json:"err,omitempty"`
}

func (h *DataHandler) GetDriverLocations(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	// Get parameters.
	id := ps.ByName("id")
	minutes, err := strconv.ParseFloat(request.URL.Query().Get(queryName), 64)
	if err != nil {
		Error(writer, err, http.StatusInternalServerError)
		return
	}

	// Get locations and send.
	d, err := h.CabService.GetDriverLocations(id, minutes)
	if err != nil {
		Error(writer, err, http.StatusInternalServerError)
	} else if d == nil {
		NotFound(writer)
	} else {
		encodeJSON(writer, &driver_location.Locations{Locations: d.Locations, Err: d.Err, ServerIP: getIP()})
	}
}

// Ensure service implements interface.
var _ driver_location.CabService = &CabService{}

// CabService represents an HTTP implementation of driver_location.CabService.
type CabService struct {
	URL *url.URL
}

//Optional
func (s *CabService) StoreLocation(id string, loc *driver_location.Location) error {
	u := *s.URL
	u.Path = "/drivers/" + url.QueryEscape(string(id)) + "/locations"

	// Encode request body.
	reqBody, err := json.Marshal(*loc)
	if err != nil {
		return err
	}

	// Create request.
	req, err := http.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	// Execute request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody patchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Err != "" {
		return driver_location.Error(respBody.Err)
	}

	return nil
}

func (s *CabService) GetDriverLocations(id string, minutes float64) (*driver_location.Locations, error) {
	u := *s.URL
	u.Path = "/drivers/" + url.QueryEscape(id) + "/locations"
	u.RawQuery = queryName + "=" + fmt.Sprintf("%f", minutes)

	// Prepare request.
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Execute request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody driver_location.Locations
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return &respBody, nil

}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}
