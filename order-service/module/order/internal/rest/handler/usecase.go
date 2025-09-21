package handler

import (
	"context"
	"order-service/module/order/entity"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type OrderUsecase interface {
	CreateOrder(ctx context.Context, params *entity.CreateOrderRequest) error
}
