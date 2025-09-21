package entity

import "time"

type Warehouse struct {
	ID        string    `json:"id"`
	ShopID    string    `json:"shop_id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WarehouseActivationRequest struct {
	WarehouseID string `json:"warehouse_id" validate:"required"`
	Active      bool   `json:"active"`
}
