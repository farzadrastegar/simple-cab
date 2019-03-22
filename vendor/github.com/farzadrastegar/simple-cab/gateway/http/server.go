package http

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/config"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

// DefaultAddr is the default bind address.
var DefaultAddr = ":4000"

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
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	configHandler := config.NewConfig(logger)
	configHandler.ReadYaml(gateway.GetConfigFilename())
	addr := configHandler.GetYamlValueStr("servers", "external", "address")
	port := configHandler.GetYamlValueStr("servers", "external", "port")
	DefaultAddr = fmt.Sprintf("%s:%s", addr, port)

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
func (c *Client) Connect() gateway.CabService {
	return &c.cabService
}
