package entity

import "time"

type WarehouseStock struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouse_id"`
	ProductID   string    `json:"product_id"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListWarehouseStockByParams struct {
	ProductIDs []string
}

type ListWarehouseStockResponse struct {
	Warehouse       []*Warehouse      `json:"warehouses"`
	WarehouseStocks []*WarehouseStock `json:"warehouse_stocks"`
	Meta            *Meta             `json:"meta"`
}

type WarehouseStockAdjustmentParams struct {
	WarehouseID string
	ProductID   string
	Stock       uint32
}

type WarehouseStockAdjustment struct {
	WarehouseID string `json:"warehouse_id" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	Stock       int    `json:"stock" validate:"required"`
}

type WarehouseStockAdjustmentRequest struct {
	WarehouseStocks []*WarehouseStockAdjustment `json:"warehouse_stocks" validate:"required,min=1,dive,required"`
}
