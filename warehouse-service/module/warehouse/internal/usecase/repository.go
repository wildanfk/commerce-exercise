package usecase

import (
	"context"
	"warehouse-service/internal/util"
	"warehouse-service/module/warehouse/entity"
)

//go:generate mockgen -destination=mock/repository.go -package=mock -source=repository.go

type WarehouseRepository interface {
	ListByIDs(ctx context.Context, ids []string) ([]*entity.Warehouse, error)
	UpdateActive(ctx context.Context, id string, active bool) error
}

type WarehouseStockRepository interface {
	Create(ctx context.Context, warehouseStock *entity.WarehouseStock, tx util.DatabaseTransaction) error
	IncreaseStock(ctx context.Context, params entity.WarehouseStockAdjustmentParams, tx util.DatabaseTransaction) (int64, error)
	DecreaseStock(ctx context.Context, params entity.WarehouseStockAdjustmentParams, tx util.DatabaseTransaction) (int64, error)
	ListByWarehouseIDsAndProductIDs(ctx context.Context, warehouseIDs []string, productIDs []string) ([]*entity.WarehouseStock, error)
	ListActiveByProductIDs(ctx context.Context, productIDs []string) ([]*entity.WarehouseStock, error)
}
