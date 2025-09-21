package fixtures

import (
	"database/sql/driver"
	"time"
	"warehouse-service/module/warehouse/entity"

	"github.com/mitchellh/copystructure"
)

var (
	WarehouseStock = &entity.WarehouseStock{
		ID:          "3",
		WarehouseID: "1",
		ProductID:   "3",
		Stock:       10,
		CreatedAt:   time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 20, 21, 22, 23, 24, time.UTC),
	}
)

func NewWarehouseStock(obj *entity.WarehouseStock) *entity.WarehouseStock {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.WarehouseStock)
	return res
}

func GetWarehouseStockRow(obj *entity.WarehouseStock) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.WarehouseID,
		obj.ProductID,
		obj.Stock,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
