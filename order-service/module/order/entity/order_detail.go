package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderDetail struct {
	ID          string          `json:"id"`
	OrderID     string          `json:"order_id"`
	ProductID   string          `json:"product_id"`
	WarehouseID string          `json:"warehouse_id"`
	Stock       int             `json:"stock"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
