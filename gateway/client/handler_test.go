package client_test

import (
	"github.com/farzadrastegar/simple-cab/gateway/client"
	"github.com/farzadrastegar/simple-cab/gateway/mock"
)

// Handler represents a test wrapper fot Client.Handler
type Handler struct {
	*client.Handler

	BusService     mock.BusService
	RequestService RequestService
}

// NewHandler creates a new test Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler: client.NewHandler(),
	}
	h.Handler.BusService = &h.BusService
	h.Handler.RequestService = &h.RequestService

	return h
}
