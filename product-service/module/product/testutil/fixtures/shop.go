package fixtures

import (
	"product-service/module/product/entity"

	"github.com/mitchellh/copystructure"
)

var (
	Shop = &entity.Shop{
		ID:   "1",
		Name: "Lorem Ipsum",
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
