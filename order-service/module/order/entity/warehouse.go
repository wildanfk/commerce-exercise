package entity

type WarehouseStock struct {
	ID            string
	WarehouseID   string
	WarehouseName string
	ShopID        string
	ProductID     string
	Stock         int
}

type WarehouseStockAdjustment struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	Stock       int    `json:"stock"`
}

type WarehouseStockAdjustmentParams struct {
	WarehouseStocks []*WarehouseStockAdjustment `json:"warehouse_stocks"`
}
