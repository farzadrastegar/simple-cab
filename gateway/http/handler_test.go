package http_test

import "github.com/farzadrastegar/simple-cab/gateway/http"

// Handler represents a test wrapper for http.Handler.
type Handler struct {
	*http.Handler

	DataHandler *DataHandler
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler:     &http.Handler{},
		DataHandler: NewDataHandler(),
	}
	h.Handler.DataHandler = h.DataHandler.DataHandler
	return h
}
