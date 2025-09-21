package fixtures

import (
	"order-service/module/order/entity"

	"github.com/mitchellh/copystructure"
)

var (
	User = &entity.User{
		ID: "6",
	}
)

func NewUser(obj *entity.User) *entity.User {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.User)
	return res
}
