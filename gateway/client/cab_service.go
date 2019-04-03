package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/farzadrastegar/simple-cab/gateway"
)

//const clientYamlFilename = "../config.yaml"

var localServerAddr = "localhost"
var localServerPort = "8082"

// Ensure service implements interface.
var _ gateway.CabService = &CabService{}

type CabService struct {
	//client *Client
	handler **Handler
}

func (s *CabService) StoreLocation(id string, data *gateway.Data) error {
	idInt, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		return err
	}

	//return s.client.Handler.BusService.Produce(gateway.DriverID(idInt), data)
	h := *s.handler
	return h.BusService.Produce(gateway.DriverID(idInt), data)
}

func (s *CabService) CheckZombieStatus(id string) (*gateway.Status, error) {
	var u url.URL
	u.Path = "/drivers/" + url.QueryEscape(id)
	u.Scheme = "HTTP"
	u.Host = fmt.Sprintf("%s:%s", localServerAddr, localServerPort)

	// Prepare request.
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Execute request.
	h := *s.handler
	resp, err := h.RequestService.ExecuteRequest(req)
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

//type GetDataResponse struct {
//	Status *gateway.Status `json:"status,omitempty"`
//	Err    string          `json:"err,omitempty"`
//}

var _ RequestService = &ReqService{}

type ReqService struct{}

func (s *ReqService) ExecuteRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
