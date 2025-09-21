package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ProductDetailWarehouse struct {
	WarehouseID      string `json:"warehouse_id"`
	WarehouseName    string `json:"warehouse_name"`
	WarehouseStockID string `json:"warehouse_stock_id"`
	WarehouseStock   int    `json:"warehouse_stock"`
}

type ProductDetailShop struct {
	ID         string                    `json:"id"`
	Name       string                    `json:"name"`
	TotalStock int                       `json:"total_stock"`
	Warehouses []*ProductDetailWarehouse `json:"warehouses"`
}

type ProductDetail struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Price      decimal.Decimal      `json:"price"`
	TotalStock int                  `json:"total_stock"`
	Shops      []*ProductDetailShop `json:"shops"`
}

type ListProductByParams struct {
	Page   int
	Offset int
	Limit  int
	IDs    []string
	Name   string
}

type ListProductResponse struct {
	Products []*Product `json:"products"`
	Meta     *ListMeta  `json:"meta"`
}

type ListProductDetailResponse struct {
	Products []*ProductDetail `json:"products"`
	Meta     *ListMeta        `json:"meta"`
}
