package usecase

import (
	"context"
	"order-service/internal/util"
	"order-service/module/order/entity"
)

//go:generate mockgen -destination=mock/repository.go -package=mock -source=repository.go

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order, tx util.DatabaseTransaction) error
	UpdateExpired(ctx context.Context, id string, tx util.DatabaseTransaction) (int64, error)
	ListByOrderExpired(ctx context.Context) ([]*entity.Order, error)
}

type OrderDetailRepository interface {
	Create(ctx context.Context, orderDetail *entity.OrderDetail, tx util.DatabaseTransaction) error
	ListByOrderID(ctx context.Context, orderID string) ([]*entity.OrderDetail, error)
}

type ProductRepository interface {
	ListByProductIDs(ctx context.Context, productIDs []string) ([]*entity.Product, error)
}

type WarehouseRepository interface {
	ActiveStock(ctx context.Context, productIDs []string) ([]*entity.WarehouseStock, error)
	AdjustmentStock(ctx context.Context, params *entity.WarehouseStockAdjustmentParams) error
}
