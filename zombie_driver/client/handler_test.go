package client_test

import (
	"github.com/farzadrastegar/simple-cab/zombie_driver/client"
)

// Handler represents a test wrapper fot Client.Handler
type Handler struct {
	*client.Handler

	RequestService RequestService
}

// NewHandler creates a new test Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler: client.NewHandler(),
	}
	h.Handler.RequestService = &h.RequestService

	return h
}
