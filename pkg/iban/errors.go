package iban

type Error struct {
	Code    string
	Message string
	Value   string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrInvalidFormat = &Error{Code: "invalid_format", Message: "IBAN does not match country format"}
)
