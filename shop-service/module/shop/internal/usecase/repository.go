package usecase

import (
	"context"
	"shop-service/internal/util/libpagination"
	"shop-service/module/shop/entity"
)

//go:generate mockgen -destination=mock/repository.go -package=mock -source=repository.go

type ShopRepository interface {
	ListByParams(ctx context.Context, params *entity.ListShopByParams) ([]*entity.Shop, *libpagination.OffsetPagination, error)
}
