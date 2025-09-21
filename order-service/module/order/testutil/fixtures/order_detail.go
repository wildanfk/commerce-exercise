package fixtures

import (
	"database/sql/driver"
	"order-service/module/order/entity"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/shopspring/decimal"
)

var (
	OrderDetail = &entity.OrderDetail{
		ID:          "11",
		OrderID:     "1",
		ProductID:   "8",
		WarehouseID: "10",
		Stock:       5,
		Price:       decimal.NewFromInt(10000),
		CreatedAt:   time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
		UpdatedAt:   time.Date(2025, 1, 10, 11, 12, 13, 14, time.UTC),
	}
)

func NewOrderDetail(obj *entity.OrderDetail) *entity.OrderDetail {
	r, err := copystructure.Copy(obj)
	if err != nil {
		return nil
	}
	res := r.(*entity.OrderDetail)
	res.Price = obj.Price

	return res
}

func GetOrderDetailRow(obj *entity.OrderDetail) []driver.Value {
	return []driver.Value{
		obj.ID,
		obj.OrderID,
		obj.ProductID,
		obj.WarehouseID,
		obj.Stock,
		obj.Price,
		obj.CreatedAt,
		obj.UpdatedAt,
	}
}
