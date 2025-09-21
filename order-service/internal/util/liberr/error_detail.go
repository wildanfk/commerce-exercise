package liberr

type ErrorDetails struct {
	Message string
	Code    string
	Field   string
}

func NewErrorDetails(message, code, field string) *ErrorDetails {
	return &ErrorDetails{
		Message: message,
		Code:    code,
		Field:   field,
	}
}

func (e *ErrorDetails) Error() string {
	return e.Message
}

func ErrorCodeEquals(err error, code string) bool {
	errDetails, ok := err.(*ErrorDetails)
	if !ok {
		return false
	}

	return errDetails.Code == code
}
