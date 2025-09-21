package usecase

import (
	"context"
	"product-service/internal/util/libpagination"
	"product-service/module/product/entity"
)

//go:generate mockgen -destination=mock/repository.go -package=mock -source=repository.go

type ProductRepository interface {
	ListByParams(ctx context.Context, params *entity.ListProductByParams) ([]*entity.Product, *libpagination.OffsetPagination, error)
}

type WarehouseRepository interface {
	ActiveStock(ctx context.Context, productIDs []string) ([]*entity.WarehouseStock, error)
}

type ShopRepository interface {
	ListByShopIDs(ctx context.Context, shopIDs []string) ([]*entity.Shop, error)
}
