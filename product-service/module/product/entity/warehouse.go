package entity

type WarehouseStock struct {
	ID            string
	WarehouseID   string
	WarehouseName string
	ShopID        string
	ProductID     string
	Stock         int
}
