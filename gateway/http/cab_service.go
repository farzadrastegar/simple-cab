package http

import (
	"bytes"
	"encoding/json"
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/config"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const ErrInvalidJSON = gateway.Error("invalid json")

//const GatewayConfigFilename = "../config.yaml"

// DataHandler represents an HTTP API handler for requests.
type DataHandler struct {
	*httprouter.Router
	Logger *log.Logger

	CabService gateway.CabService
}

// NewDataHandler returns a new instance of DataHandler.
func NewDataHandler() *DataHandler {
	h := &DataHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}

	//create a config handler
	configHandler := config.NewConfig(h.Logger)

	//read yaml config
	configHandler.ReadYaml(gateway.GetConfigFilename())

	//setup routes
	h.SetupRoutes(configHandler)

	return h
}

func (h *DataHandler) SetupRoutes(conf *config.Handlers) {
	//check method types in yaml
	if !strings.EqualFold(conf.GetYamlValueStr("urls", "driverLocations", "method"), http.MethodPatch) {
		h.Logger.Fatalf("ERROR: wrong method type in yaml config file")
	}
	if !strings.EqualFold(conf.GetYamlValueStr("urls", "zombieDriver", "method"), http.MethodGet) {
		h.Logger.Fatalf("ERROR: wrong method type in yaml config file")
	}

	//set route of first endpoint
	h.PATCH(conf.GetYamlValueStr("urls", "driverLocations", "path"), h.StoreLocation)

	//set route of second endpoint
	h.GET(conf.GetYamlValueStr("urls", "zombieDriver", "path"), h.CheckDriverStatus)
}

func (h *DataHandler) StoreLocation(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	// Decode request.
	var req gateway.Data

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		Error(writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	} else {
		//get URL parameters
		driverId := ps.ByName("id")

		//send request body via NSQ
		switch err := h.CabService.StoreLocation(driverId, &req); err {
		case nil:
			encodeJSON(writer, &patchResponse{}, h.Logger)
		default:
			Error(writer, gateway.ErrInternal, http.StatusBadRequest, h.Logger)
		}
	}
}

type patchResponse struct {
	Err string `json:"err,omitempty"`
}

func (h *DataHandler) CheckDriverStatus(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	// Status by ID.
	d, err := h.CabService.CheckZombieStatus(id)
	if err != nil {
		Error(writer, err, http.StatusInternalServerError, h.Logger)
	} else if d == nil {
		NotFound(writer)
	} else {
		encodeJSON(writer, &getResponse{ID: d.ID, Zombie: d.Zombie}, h.Logger)
	}
}

type getResponse struct {
	ID     gateway.DriverID `json:"id"`
	Zombie bool             `json:"zombie"`
	Err    string           `json:"err,omitempty"`
}

// Ensure service implements interface.
var _ gateway.CabService = &CabService{}

// CabService represents an HTTP implementation of gateway.CabService.
type CabService struct {
	URL *url.URL
}

func (s *CabService) StoreLocation(id string, data *gateway.Data) error {
	u := *s.URL
	u.Path = "/drivers/" + url.QueryEscape(string(id)) + "/locations"

	// Encode request body.
	reqBody, err := json.Marshal(*data)
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
		return gateway.Error(respBody.Err)
	}

	return nil
}

func (s *CabService) CheckZombieStatus(id string) (*gateway.Status, error) {
	u := *s.URL
	u.Path = "/drivers/" + url.QueryEscape(id)

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
	var respBody gateway.Status
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}
