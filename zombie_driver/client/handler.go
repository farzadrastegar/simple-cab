package client

type Handler struct {
	RequestService RequestService
}

func NewHandler() *Handler {
	return &Handler{
		RequestService: &ReqService{},
	}
}
