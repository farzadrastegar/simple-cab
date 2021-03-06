package client

import (
	"fmt"
	"net/http"

	"github.com/farzadrastegar/simple-cab/zombie_driver"

	"github.com/spf13/viper"
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
	// Read CheckZombieStatus service's config.
	localServerAddr = viper.GetString("servers.driver_location.address")
	localServerPort = viper.GetString("servers.driver_location.port")
	distanceFloat = viper.GetFloat64("urls.zombieStatus.definition.distance")
	distance = fmt.Sprintf("%f", distanceFloat)
	durationFloat := viper.GetFloat64("urls.zombieStatus.definition.duration")
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
