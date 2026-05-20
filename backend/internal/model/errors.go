package model

// Err is a domain error carrying an HTTP status code and a client-safe
// message. Predefined Err values live in the constants package.
type Err struct {
	Code    int
	Message string
	Details []string
}

// Error implements the error interface.
func (e *Err) Error() string {
	return e.Message
}

// WithDetails returns a copy of the error with extra detail strings attached.
func (e *Err) WithDetails(details ...string) *Err {
	return &Err{
		Code:    e.Code,
		Message: e.Message,
		Details: append(append([]string{}, e.Details...), details...),
	}
}

// NewErr builds a custom domain error.
func NewErr(code int, message string, details ...string) *Err {
	return &Err{Code: code, Message: message, Details: details}
}
