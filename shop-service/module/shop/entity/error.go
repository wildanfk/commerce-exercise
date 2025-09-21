package entity

import (
	"shop-service/internal/util/liberr"
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
	ErrorCodeForbidden        = "FORBIDDEN"
	ErrorCodeInvalidParameter = "PARAMETER_INVALID"
	ErrorCodeShopNotFound     = "SHOP_NOT-FOUND"
)

var (
	ErrorForbidden        = liberr.NewErrorDetails("Forbidden", ErrorCodeForbidden, "")
	ErrorInvalidParameter = liberr.NewErrorDetails("Invalid parameter", ErrorCodeInvalidParameter, "")
	ErrorUserNotFound     = liberr.NewErrorDetails("Shop Not Found", ErrorCodeShopNotFound, "")
)
