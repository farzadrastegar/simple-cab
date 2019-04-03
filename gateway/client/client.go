package client

import (
	"net/http"

	"github.com/farzadrastegar/simple-cab/gateway"

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
	// Read CheckZombieStatus service's address and port.
	localServerAddr = viper.GetString("servers.zombie_driver.address")
	localServerPort = viper.GetString("servers.zombie_driver.port")

	c := &Client{
		Handler: NewHandler(),
	}

	//c.cabService.client = c
	c.cabService.handler = &c.Handler
	return c
}

// Connect returns the cabservice from client.
func (c *Client) Connect() gateway.CabService {
	return &c.cabService
}
