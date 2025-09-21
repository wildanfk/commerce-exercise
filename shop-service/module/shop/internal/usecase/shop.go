package usecase

import (
	"context"
	"shop-service/internal/util/libpagination"
	"shop-service/module/shop/entity"
)

type ShopUsecaseRepos struct {
	ShopRepo ShopRepository
}

type ShopUsecase struct {
	repos *ShopUsecaseRepos
}

func NewShopUsecase(repos *ShopUsecaseRepos) *ShopUsecase {
	return &ShopUsecase{
		repos: repos,
	}
}

func (s *ShopUsecase) ListByParams(ctx context.Context, params *entity.ListShopByParams) ([]*entity.Shop, *libpagination.OffsetPagination, error) {
	params.Offset = libpagination.Offset(params.Page, params.Limit)

	return s.repos.ShopRepo.ListByParams(ctx, params)
}
