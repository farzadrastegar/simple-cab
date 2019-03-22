package client

import (
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/config"
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
	conf.ReadYaml(gateway.GetConfigFilename())
	localServerAddr = conf.GetYamlValueStr("servers", "internal", "address")
	localServerPort = conf.GetYamlValueStr("servers", "internal", "port")

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
