package usecase

import (
	"context"
	"database/sql"
	"math"
	"warehouse-service/internal/util"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/internal/util/libvalidate"
	"warehouse-service/module/warehouse/entity"
)

type WarehouseStockUsecaseRepos struct {
	DatabaseTransactionHandler util.DatabaseTransactionHandler
	WarehouseRepo              WarehouseRepository
	WarehouseStockRepo         WarehouseStockRepository
}

type WarehouseStockUsecase struct {
	repos *WarehouseStockUsecaseRepos
}

func NewWarehouseStockUsecase(repos *WarehouseStockUsecaseRepos) *WarehouseStockUsecase {
	return &WarehouseStockUsecase{
		repos: repos,
	}
}

func (ws *WarehouseStockUsecase) ActiveStock(ctx context.Context, params *entity.ListWarehouseStockByParams) ([]*entity.Warehouse, []*entity.WarehouseStock, error) {
	warehouseStocks, err := ws.repos.WarehouseStockRepo.ListActiveByProductIDs(ctx, params.ProductIDs)
	if err != nil {
		return nil, nil, liberr.ResolveError(err)
	}

	warehouseIDsMap := make(map[string]struct{})
	warehouseIDs := []string{}
	for _, ws := range warehouseStocks {
		if _, exists := warehouseIDsMap[ws.WarehouseID]; !exists {
			warehouseIDsMap[ws.WarehouseID] = struct{}{}
			warehouseIDs = append(warehouseIDs, ws.WarehouseID)
		}
	}

	warehouses := []*entity.Warehouse{}
	if len(warehouseIDs) > 0 {
		warehouses, err = ws.repos.WarehouseRepo.ListByIDs(ctx, warehouseIDs)
		if err != nil {
			return nil, nil, liberr.ResolveError(err)
		}
	}

	return warehouses, warehouseStocks, nil
}

func (ws *WarehouseStockUsecase) AdjustmentStock(ctx context.Context, params *entity.WarehouseStockAdjustmentRequest) error {
	// Validation struct
	if err := libvalidate.Validator().Struct(params); err != nil {
		return libvalidate.ResolveError(err, entity.ErrorCodeInvalidBodyJSON)
	}

	// Validation Stock Adjustment
	if err := ws.stockAdjustmentValidation(ctx, params.WarehouseStocks); err != nil {
		return liberr.ResolveError(err)
	}

	// Process Stock Adjustment
	return ws.stockAdjustment(ctx, params.WarehouseStocks)
}

func (ws *WarehouseStockUsecase) TransferStock(ctx context.Context, params *entity.WarehouseStockTransferRequest) error {
	// Validation struct
	if err := libvalidate.Validator().Struct(params); err != nil {
		return libvalidate.ResolveError(err, entity.ErrorCodeInvalidBodyJSON)
	}

	// Build stock adjustment
	stockAdjustments := []*entity.WarehouseStockAdjustment{}

	for _, p := range params.Products {
		stockAbs := int(math.Abs(float64(p.Stock)))

		// Decrese from origin warehouse
		stockAdjustments = append(stockAdjustments, &entity.WarehouseStockAdjustment{
			WarehouseID: params.OriginalWarehouseID,
			ProductID:   p.ProductID,
			Stock:       -1 * stockAbs,
		})

		// Increase to destination warehouse
		stockAdjustments = append(stockAdjustments, &entity.WarehouseStockAdjustment{
			WarehouseID: params.DestinationWarehouseID,
			ProductID:   p.ProductID,
			Stock:       1 * stockAbs,
		})
	}

	// Validation Stock Adjustment
	if err := ws.stockAdjustmentValidation(ctx, stockAdjustments); err != nil {
		return liberr.ResolveError(err)
	}

	// Process Stock Adjustment
	return ws.stockAdjustment(ctx, stockAdjustments)
}

func (ws *WarehouseStockUsecase) stockAdjustmentValidation(ctx context.Context, stockAdjustments []*entity.WarehouseStockAdjustment) error {
	// Get Unique Warehouse ID & Warehouse Stock ID
	warehouseIDsMap := make(map[string]struct{})
	warehouseIDs := []string{}

	productIDsMap := make(map[string]struct{})
	productIDs := []string{}

	for _, sa := range stockAdjustments {
		if _, exists := warehouseIDsMap[sa.WarehouseID]; !exists {
			warehouseIDsMap[sa.WarehouseID] = struct{}{}
			warehouseIDs = append(warehouseIDs, sa.WarehouseID)
		}

		if _, exists := productIDsMap[sa.ProductID]; !exists {
			productIDsMap[sa.ProductID] = struct{}{}
			productIDs = append(productIDs, sa.ProductID)
		}
	}

	// Validate Warehouse
	warehouses, err := ws.repos.WarehouseRepo.ListByIDs(ctx, warehouseIDs)
	if len(warehouseIDs) != len(warehouses) {
		return liberr.ResolveError(entity.ErrorWarehouseNotFound)
	}

	// Validate Warehouse Stock
	warehouseStocks, err := ws.repos.WarehouseStockRepo.ListByWarehouseIDsAndProductIDs(ctx, warehouseIDs, productIDs)
	if err != nil {
		return liberr.ResolveError(err)
	}

	warehouseProductStockMap := make(map[string]map[string]int, 0)
	for _, s := range warehouseStocks {
		if _, exists := warehouseProductStockMap[s.WarehouseID]; !exists {
			warehouseProductStockMap[s.WarehouseID] = map[string]int{}
		}
		warehouseProductStockMap[s.WarehouseID][s.ProductID] = s.Stock
	}

	for _, sa := range stockAdjustments {
		if currentStock, exists := warehouseProductStockMap[sa.WarehouseID][sa.ProductID]; !exists {
			return liberr.ResolveError(entity.ErrorWarehouseStockNotFound)
		} else if currentStock+sa.Stock < 0 {
			return liberr.ResolveError(entity.ErrorWarehouseStockAdjustmentOutOfStock)
		}
	}

	return nil
}

func (ws *WarehouseStockUsecase) stockAdjustment(ctx context.Context, stockAdjustments []*entity.WarehouseStockAdjustment) error {
	tx, err := ws.repos.DatabaseTransactionHandler.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return liberr.ResolveError(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback() //nolint
		}
	}()

	var affected int64
	for _, sa := range stockAdjustments {
		stockAbs := uint32(math.Abs(float64(sa.Stock)))

		if sa.Stock >= 0 {
			affected, err = ws.repos.WarehouseStockRepo.IncreaseStock(ctx, entity.WarehouseStockAdjustmentParams{
				WarehouseID: sa.WarehouseID,
				ProductID:   sa.ProductID,
				Stock:       stockAbs,
			}, tx)
			if err != nil {
				return liberr.ResolveError(err)
			}
			if affected <= 0 {
				err = liberr.ResolveError(entity.ErrorWarehouseStockAdjustmentFailed)
				return err
			}
		} else {
			affected, err = ws.repos.WarehouseStockRepo.DecreaseStock(ctx, entity.WarehouseStockAdjustmentParams{
				WarehouseID: sa.WarehouseID,
				ProductID:   sa.ProductID,
				Stock:       stockAbs,
			}, tx)
			if err != nil {
				return liberr.ResolveError(err)
			}
			if affected <= 0 {
				err = liberr.ResolveError(entity.ErrorWarehouseStockAdjustmentFailed)
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return liberr.ResolveError(err)
	}

	return nil
}
