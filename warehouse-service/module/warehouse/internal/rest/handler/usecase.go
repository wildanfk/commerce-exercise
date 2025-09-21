package handler

import (
	"context"
	"warehouse-service/module/warehouse/entity"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type WarehouseUsecase interface {
	WarehouseActivation(ctx context.Context, params *entity.WarehouseActivationRequest) error
}

type WarehouseStockUsecase interface {
	ActiveStock(ctx context.Context, params *entity.ListWarehouseStockByParams) ([]*entity.Warehouse, []*entity.WarehouseStock, error)
	AdjustmentStock(ctx context.Context, params *entity.WarehouseStockAdjustmentRequest) error
	TransferStock(ctx context.Context, params *entity.WarehouseStockTransferRequest) error
}
