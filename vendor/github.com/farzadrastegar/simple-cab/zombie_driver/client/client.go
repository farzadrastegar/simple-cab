package client

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"github.com/farzadrastegar/simple-cab/zombie_driver/config"
	"log"
	"net/http"
	"os"
)

type RequestService interface {
	ExecuteRequest(req *http.Request) (*http.Response, error)
}

// Client represents a client to connect to the HTTP server.
type Client struct {
	cabService CabService

	Handler *Handler
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	// Read CheckZombieStatus service's address and port.
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	conf := config.NewConfig(logger)
	conf.ReadYaml(zombie_driver.GetConfigFilename())
	localServerAddr = conf.GetYamlValueStr("servers", "driver_location", "address")
	localServerPort = conf.GetYamlValueStr("servers", "driver_location", "port")
	distanceFloat, _ = conf.GetYamlValue("urls", "zombieStatus", "definition", "distance").Float()
	distance = fmt.Sprintf("%f", distanceFloat)
	durationFloat, _ := conf.GetYamlValue("urls", "zombieStatus", "definition", "duration").Float()
	duration = fmt.Sprintf("%f", durationFloat)

	c := &Client{
		Handler: NewHandler(),
	}

	//c.cabService.client = c
	c.cabService.handler = &c.Handler
	return c
}

// Connect returns the cabservice from client.
func (c *Client) Connect() zombie_driver.CabService {
	return &c.cabService
}
