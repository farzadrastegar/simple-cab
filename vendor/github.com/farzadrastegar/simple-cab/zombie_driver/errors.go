package zombie_driver

// General errors.
const (
	ErrUnauthorized = Error("unauthorized")
	ErrInternal     = Error("internal error")
)

// Data errors.
const (
	ErrInvalidData = Error("invalid input data")
)

// Error represents a zombie_driver error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
