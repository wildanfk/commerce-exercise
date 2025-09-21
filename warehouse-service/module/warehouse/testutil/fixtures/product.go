package fixtures

import (
	"warehouse-service/module/warehouse/entity"

	"github.com/mitchellh/copystructure"
)

var (
	Product = &entity.Product{
		ID:   "3",
		Name: "Lorem Ipsum Product",
	}
)

func NewProduct(obj *entity.Product) *entity.Product {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.Product)
	return res
}
