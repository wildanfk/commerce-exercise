package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderState int

const (
	OrderStateUnspecified OrderState = iota
	OrderStateCreated
	OrderStateExpired
)

type Order struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	ShopID     string          `json:"shop_id"`
	State      OrderState      `json:"state"`
	TotalStock int             `json:"total_stock"`
	TotalPrice decimal.Decimal `json:"total_price"`
	ExpiredAt  time.Time       `json:"expired_at"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type CreateOrderProduct struct {
	ProductID   string `json:"id" validate:"required"`
	WarehouseID string `json:"warehouse_id" validate:"required"`
	Stock       int    `json:"stock" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	ShopID   string                `json:"shop_id" validate:"required"`
	Products []*CreateOrderProduct `json:"products" validate:"required,min=1,dive,required"`
	User     *User
}
