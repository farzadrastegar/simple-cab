package http

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/farzadrastegar/simple-cab/zombie_driver"

	"github.com/spf13/viper"
)

// DefaultAddr is the default bind address.
var DefaultAddr = ":4001"

// Server represents an HTTP server.
type Server struct {
	ln net.Listener

	// Handler to serve.
	Handler *Handler

	// Bind address to open.
	Addr string
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	// Read server address and port from config.yaml.
	port := viper.GetString("servers.zombie_driver.port")
	if port != "" {
		addr := viper.GetString("servers.zombie_driver.address")
		DefaultAddr = fmt.Sprintf("%s:%s", addr, port)
	}

	return &Server{
		Addr: DefaultAddr,
	}
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	// Start HTTP server.
	go func() { http.Serve(s.ln, s.Handler) }()

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port returns the port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

// Client represents a client to connect to the HTTP server.
type Client struct {
	URL        url.URL
	cabService CabService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{}
	c.cabService.URL = &c.URL
	return c
}

// Connect returns the service for managing data.
func (c *Client) Connect() zombie_driver.CabService {
	return &c.cabService
}
