package fixtures

import (
	"database/sql/driver"
	"time"
	"warehouse-service/module/warehouse/entity"

	"github.com/mitchellh/copystructure"
)

var (
	Warehouse = &entity.Warehouse{
		ID:        "1",
		ShopID:    "11",
		Name:      "Lorem Ipsum Warehouse",
		Active:    true,
		CreatedAt: time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt: time.Date(2025, 2, 20, 21, 22, 23, 24, time.UTC),
	}
)

func NewWarehouse(obj *entity.Warehouse) *entity.Warehouse {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.Warehouse)
	return res
}

func GetWarehouseRow(obj *entity.Warehouse) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.ShopID,
		obj.Name,
		obj.Active,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
