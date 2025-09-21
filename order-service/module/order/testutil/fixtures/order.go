package fixtures

import (
	"database/sql/driver"
	"order-service/module/order/entity"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/shopspring/decimal"
)

var (
	Order = &entity.Order{
		ID:         "1",
		UserID:     "2",
		ShopID:     "3",
		State:      entity.OrderStateCreated,
		TotalStock: 5,
		TotalPrice: decimal.NewFromInt(50000),
		ExpiredAt:  time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		CreatedAt:  time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt:  time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
	}
)

func NewOrder(obj *entity.Order) *entity.Order {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.Order)
	res.TotalPrice = obj.TotalPrice

	return res
}

func GetOrderRow(obj *entity.Order) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.UserID,
		obj.ShopID,
		obj.State,
		obj.TotalStock,
		obj.TotalPrice,
		obj.ExpiredAt,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
