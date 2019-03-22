package client

import "github.com/farzadrastegar/simple-cab/gateway"

type Handler struct {
	BusService     gateway.BusService
	RequestService RequestService
}

func NewHandler() *Handler {
	return &Handler{
		RequestService: &ReqService{},
	}
}
