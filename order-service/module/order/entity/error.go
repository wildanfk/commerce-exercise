package entity

import (
	"order-service/internal/util/liberr"
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
	ErrorCodeForbidden                = "FORBIDDEN"
	ErrorCodeTokenNotFound            = "TOKEN_NOT-FOUND"
	ErrorCodeTokenExpired             = "TOKEN_EXPIRED"
	ErrorCodeTokenInvalid             = "TOKEN_INVALID"
	ErrorCodeTokenInvalidBarer        = "TOKEN_INVALID_BEARER"
	ErrorCodeInvalidBodyJSON          = "BODY-JSON_INVALID"
	ErrorCodeInvalidParameter         = "PARAMETER_INVALID"
	ErrorCodeOrderNotFound            = "ORDER_NOT-FOUND"
	ErrorCodeProductNotFound          = "ORDER-PRODUCT_NOT-FOUND"
	ErrorCodeProductStockNotFound     = "ORDER-PRODUCT-STOCK_NOT-FOUND"
	ErrorCodeProductOutOfStock        = "ORDER-PRODUCT_OUT-OF-STOCK"
	ErrorCodeProductInsufficientStock = "ORDER-PRODUCT_INSUFFICIENT-STOCK"
	ErrorCodeProductConflicted        = "ORDER-PRODUCT_CONFLICTED"
	ErrorCodeProductMultiShop         = "ORDER-PRODUCT_MULTI-SHOP"
)

var (
	ErrorForbidden                = liberr.NewErrorDetails("Forbidden", ErrorCodeForbidden, "")
	ErrorTokenNotFound            = liberr.NewErrorDetails("Token Not Found", ErrorCodeTokenNotFound, "")
	ErrorTokenExpired             = liberr.NewErrorDetails("Token Expired", ErrorCodeTokenExpired, "")
	ErrorTokenInvalid             = liberr.NewErrorDetails("Token Invalid", ErrorCodeTokenInvalid, "")
	ErrorTokenInvalidBearer       = liberr.NewErrorDetails("Token Invalid Due Bearer", ErrorCodeTokenInvalidBarer, "")
	ErrorInvalidBodyJSON          = liberr.NewErrorDetails("Invalid body JSON", ErrorCodeInvalidBodyJSON, "")
	ErrorInvalidParameter         = liberr.NewErrorDetails("Invalid parameter", ErrorCodeInvalidParameter, "")
	ErrorUserNotFound             = liberr.NewErrorDetails("Order Not Found", ErrorCodeOrderNotFound, "")
	ErrorProductNotFound          = liberr.NewErrorDetails("Order Product Not Found", ErrorCodeProductNotFound, "")
	ErrorProductStockNotFound     = liberr.NewErrorDetails("Order Product Stock Not Found", ErrorCodeProductStockNotFound, "")
	ErrorProductOutOfStock        = liberr.NewErrorDetails("Order Product Out Of Stock", ErrorCodeProductOutOfStock, "")
	ErrorProductInsufficientStock = liberr.NewErrorDetails("Order Product Insufficient Stock", ErrorCodeProductInsufficientStock, "")
	ErrorProductConflicted        = liberr.NewErrorDetails("Order Product Conflicted or Out Of Stock", ErrorCodeProductConflicted, "")
	ErrorProductMultiShop         = liberr.NewErrorDetails("Order Product From Multi Shop", ErrorCodeProductMultiShop, "")
)
