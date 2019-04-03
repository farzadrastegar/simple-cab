package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/farzadrastegar/simple-cab/zombie_driver"

	"github.com/julienschmidt/httprouter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ErrInvalidJSON = zombie_driver.Error("invalid json")

//const ZombieConfigFilename = "../config.yaml"

// DataHandler represents an HTTP API handler for requests.
type DataHandler struct {
	*httprouter.Router

	CabService zombie_driver.CabService
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
	if !strings.EqualFold(viper.GetString("urls.zombieStatus.method"), http.MethodGet) {
		logger.Fatalf("ERROR: wrong method type in yaml config file")
	}

	//setup routes for its own services
	h.GET(viper.GetString("urls.zombieStatus.path"), h.CheckZombieStatus)
}

func (h *DataHandler) CheckZombieStatus(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	// Status by ID.
	d, err := h.CabService.CheckZombieStatus(id)
	if err != nil {
		Error(writer, err, http.StatusInternalServerError)
	} else if d == nil {
		NotFound(writer)
	} else {
		encodeJSON(writer, &zombie_driver.Status{ID: d.ID, Zombie: d.Zombie})
	}
}

// Ensure service implements interface.
var _ zombie_driver.CabService = &CabService{}

// CabService represents an HTTP implementation of zombie_driver.CabService.
type CabService struct {
	URL *url.URL
}

func (s *CabService) CheckZombieStatus(id string) (*zombie_driver.Status, error) {
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
	var respBody zombie_driver.Status
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return &respBody, nil

}
