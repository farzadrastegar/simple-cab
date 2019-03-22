package gateway

// General errors.
const (
	ErrUnauthorized = Error("unauthorized")
	ErrInternal     = Error("internal error")
)

// Data errors.
const (
	ErrInvalidData = Error("invalid input data")
)

// Error represents a gateway error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
