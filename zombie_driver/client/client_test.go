package client_test

import (
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"github.com/farzadrastegar/simple-cab/zombie_driver/client"
)

// Client represents a test wrapper for client.Client.
type Client struct {
	*client.Client

	cabService zombie_driver.CabService
	Handler    *Handler
}

// NewClient returns a new instance of test Client.
func NewClient() *Client {
	c := &Client{
		Client:  client.NewClient(),
		Handler: NewHandler(),
	}
	c.Client.Handler = client.NewHandler()
	c.cabService = c.Client.Connect()
	c.Client.Handler.RequestService = &c.Handler.RequestService
	return c
}

// Connect returns the cabservice from client.
func (c *Client) Connect() zombie_driver.CabService {
	return c.cabService
}
