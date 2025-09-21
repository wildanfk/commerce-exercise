package fixtures

import (
	"database/sql/driver"
	"shop-service/module/shop/entity"
	"time"

	"github.com/mitchellh/copystructure"
)

var (
	Shop = &entity.Shop{
		ID:        "2",
		Name:      "Lorem Ipsum",
		CreatedAt: time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt: time.Date(2025, 2, 20, 21, 22, 23, 24, time.UTC),
	}
)

func NewShop(obj *entity.Shop) *entity.Shop {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.Shop)
	return res
}

func GetShopRow(obj *entity.Shop) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.Name,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
