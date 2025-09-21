package handler

import (
	"context"
	"product-service/internal/util/libpagination"
	"product-service/module/product/entity"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type ProductUsecase interface {
	CheckProduct(ctx context.Context, params *entity.ListProductByParams) ([]*entity.Product, *libpagination.OffsetPagination, error)
	ListProduct(ctx context.Context, params *entity.ListProductByParams) ([]*entity.ProductDetail, *libpagination.OffsetPagination, error)
}
