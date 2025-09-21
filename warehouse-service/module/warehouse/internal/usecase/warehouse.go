package usecase

import (
	"context"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/internal/util/libvalidate"
	"warehouse-service/module/warehouse/entity"
)

type WarehouseUsecaseRepos struct {
	WarehouseRepo WarehouseRepository
}

type WarehouseUsecase struct {
	repos *WarehouseUsecaseRepos
}

func NewWarehouseUsecase(repos *WarehouseUsecaseRepos) *WarehouseUsecase {
	return &WarehouseUsecase{
		repos: repos,
	}
}

func (w *WarehouseUsecase) WarehouseActivation(ctx context.Context, params *entity.WarehouseActivationRequest) error {
	if err := libvalidate.Validator().Struct(params); err != nil {
		return libvalidate.ResolveError(err, entity.ErrorCodeInvalidBodyJSON)
	}

	warehouses, err := w.repos.WarehouseRepo.ListByIDs(ctx, []string{params.WarehouseID})
	if err != nil {
		return liberr.ResolveError(err)
	}

	if len(warehouses) == 0 {
		return liberr.ResolveError(entity.ErrorWarehouseNotFound)
	}

	return w.repos.WarehouseRepo.UpdateActive(ctx, params.WarehouseID, params.Active)
}
