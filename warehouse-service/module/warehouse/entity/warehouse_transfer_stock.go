package entity

type WarehouseProductStockTransferProduct struct {
	ProductID string `json:"product_id" validate:"required"`
	Stock     int    `json:"stock" validate:"required,gt=0"`
}

type WarehouseStockTransferRequest struct {
	OriginalWarehouseID    string                                  `json:"original_warehouse_id" validate:"required"`
	DestinationWarehouseID string                                  `json:"destination_warehouse_id" validate:"required"`
	Products               []*WarehouseProductStockTransferProduct `json:"products" validate:"required,min=1,dive,required"`
}
