package entity

import (
	"user-service/internal/util/liberr"
)

type Error struct {
	ErrorMessage string `json:"message"`
	ErrorCode    string `json:"code"`
	ErrorField   string `json:"field"`
}

type ErrorResponse struct {
	Errors []*Error `json:"errors"`
	Meta   *Meta    `json:"meta"`
}

const (
	ErrorCodeInvalidBodyJSON = "BODY-JSON_INVALID"
	ErrorCodeUserNotFound    = "USER_NOT-FOUND"
	ErrorCodeAuthInvalid     = "AUTH_INVALID"
)

var (
	ErrorInvalidBodyJSON = liberr.NewErrorDetails("Invalid body JSON", ErrorCodeInvalidBodyJSON, "")
	ErrorUserNotFound    = liberr.NewErrorDetails("User Not Found", ErrorCodeUserNotFound, "")
	ErrorAuthInvalid     = liberr.NewErrorDetails("Authentication invalid due wrong phone/email or password", ErrorCodeAuthInvalid, "")
)
