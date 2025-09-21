package fixtures

import (
	"database/sql/driver"
	"time"
	"user-service/module/auth/entity"

	"github.com/mitchellh/copystructure"
)

var (
	User = &entity.User{
		ID:        "1",
		Name:      "Jhon Doe",
		Email:     "jhon.doe@test.com",
		Phone:     "+6281234567890",
		Password:  "937e8d5fbb48bd4949536cd65b8d35c426b80d2f830c5c308e2cdec422ae2244", // test1234
		CreatedAt: time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt: time.Date(2025, 2, 20, 21, 22, 23, 24, time.UTC),
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

func GetUserRow(obj *entity.User) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.Name,
		obj.Email,
		obj.Phone,
		obj.Password,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
