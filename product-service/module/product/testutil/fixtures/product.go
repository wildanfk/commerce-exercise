package fixtures

import (
	"database/sql/driver"
	"product-service/module/product/entity"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/shopspring/decimal"
)

var (
	Product = &entity.Product{
		ID:        "2",
		Name:      "Lorem Ipsum Product",
		Price:     decimal.NewFromInt(10000),
		CreatedAt: time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt: time.Date(2025, 2, 20, 21, 22, 23, 24, time.UTC),
	}
)

func NewProduct(obj *entity.Product) *entity.Product {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.Product)
	res.Price = obj.Price

	return res
}

func GetProductRow(obj *entity.Product) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.Name,
		obj.Price,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
