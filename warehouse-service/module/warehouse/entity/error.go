package entity

import (
	"warehouse-service/internal/util/liberr"
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
	ErrorCodeForbidden                          = "FORBIDDEN"
	ErrorCodeInvalidBodyJSON                    = "BODY-JSON_INVALID"
	ErrorCodeInvalidParameter                   = "PARAMETER_INVALID"
	ErrorCodeWarehouseNotFound                  = "WAREHOUSE_NOT-FOUND"
	ErrorCodeWarehouseStockNotFound             = "WAREHOUSE-STOCk_NOT-FOUND"
	ErrorCodeWarehouseStockDuplicated           = "WAREHOUSE-STOCK_DUPLICATED"
	ErrorCodeWarehouseStockAdjustmentFailed     = "WAREHOUSE-STOCK_ADJUSTMENT-FAILED"
	ErrorCodeWarehouseStockAdjustmentOutOfStock = "WAREHOUSE-STOCK_ADJUSTMENT-OUT-OF-STOCK"
)

var (
	ErrorForbidden                          = liberr.NewErrorDetails("Forbidden", ErrorCodeForbidden, "")
	ErrorInvalidBodyJSON                    = liberr.NewErrorDetails("Invalid body JSON", ErrorCodeInvalidBodyJSON, "")
	ErrorInvalidParameter                   = liberr.NewErrorDetails("Invalid parameter", ErrorCodeInvalidParameter, "")
	ErrorWarehouseNotFound                  = liberr.NewErrorDetails("Warehouse Not Found", ErrorCodeWarehouseNotFound, "")
	ErrorWarehouseStockNotFound             = liberr.NewErrorDetails("Warehouse Stock Not Found", ErrorCodeWarehouseStockNotFound, "")
	ErrorWarehouseStockDuplicated           = liberr.NewErrorDetails("Warehouse Stock Already Exists", ErrorCodeWarehouseStockDuplicated, "")
	ErrorWarehouseStockAdjustmentFailed     = liberr.NewErrorDetails("Failed to Adjust Stock Due Race Condition Happened", ErrorCodeWarehouseStockAdjustmentFailed, "")
	ErrorWarehouseStockAdjustmentOutOfStock = liberr.NewErrorDetails("Failed to Adjust Stock Due Out of Stock", ErrorCodeWarehouseStockAdjustmentOutOfStock, "")
)
