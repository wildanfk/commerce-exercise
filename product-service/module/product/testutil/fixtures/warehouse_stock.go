package fixtures

import (
	"product-service/module/product/entity"

	"github.com/mitchellh/copystructure"
)

var (
	WarehouseStock = &entity.WarehouseStock{
		ID:          "10",
		WarehouseID: "3",
		ShopID:      "1",
		ProductID:   "2",
		Stock:       4,
	}
)

func NewWarehouseStock(obj *entity.Shop) *entity.WarehouseStock {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.WarehouseStock)
	return res
}
