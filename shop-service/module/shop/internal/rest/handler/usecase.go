package handler

import (
	"context"
	"shop-service/internal/util/libpagination"
	"shop-service/module/shop/entity"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type ShopUsecase interface {
	ListByParams(ctx context.Context, params *entity.ListShopByParams) ([]*entity.Shop, *libpagination.OffsetPagination, error)
}
