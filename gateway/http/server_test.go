package http_test

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/http"
	"github.com/farzadrastegar/simple-cab/gateway/mock"
	"io"
	"net/url"
	"os"
	"testing"
)

// Server represents a test wrapper for http.Server.
type Server struct {
	*http.Server

	Handler *Handler
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	gateway.SetConfigFilename("../cmd/config.yaml")

	s := &Server{
		Server:  http.NewServer(),
		Handler: NewHandler(),
	}
	s.Server.Handler = s.Handler.Handler

	// Use random port.
	s.Addr = ":0"

	return s
}

// MustOpenServerClient returns a running server and associated client. Panic on error.
func MustOpenTestServerHttpClient() (*Server, *http.Client) {
	// Create and open test server.
	s := NewServer()
	if err := s.Open(); err != nil {
		panic(err)
	}

	// Create a client pointing to the server.
	c := http.NewClient()
	c.URL = url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", s.Port())}
	return s, c
}

// MustOpenServerClient returns a running server and associated client. Panic on error.
func MustOpenServerClient() (*Server, *Client) {
	// Create and open test server.
	s := NewServer()
	if err := s.Open(); err != nil {
		panic(err)
	}

	// Create a client pointing to the server.
	c := NewClient()
	c.Client.URL = url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", s.Port())}
	//fmt.Println(c.Client.URL)
	return s, c
}

// VerboseWriter returns a multi-writer to STDERR and w if the "-v" flag is set.
func VerboseWriter(w io.Writer) io.Writer {
	if testing.Verbose() {
		return io.MultiWriter(w, os.Stderr)
	}
	return w
}

// Client represents a test wrapper for http.Client.
type Client struct {
	*http.Client

	URL        *url.URL
	cabService mock.CabService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{Client: http.NewClient()}
	c.URL = &c.Client.URL
	return c
}

// Connect returns the service for managing data.
func (c *Client) Connect() gateway.CabService {
	return &c.cabService
}
