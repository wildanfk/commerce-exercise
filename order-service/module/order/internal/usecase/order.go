package usecase

import (
	"context"
	"database/sql"
	"order-service/internal/util"
	"order-service/internal/util/liberr"
	"order-service/internal/util/libvalidate"
	"order-service/module/order/entity"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type OrderUsecaseRepos struct {
	DatabaseTransactionHandler util.DatabaseTransactionHandler
	OrderRepo                  OrderRepository
	OrderDetailRepo            OrderDetailRepository
	ProductRepo                ProductRepository
	WarehouseRepo              WarehouseRepository
}

type OrderUsecaseConfig struct {
	OrderExpirationTimeSecond int
}

type OrderUsecase struct {
	repos   *OrderUsecaseRepos
	configs *OrderUsecaseConfig
	logger  *zap.Logger
}

func NewOrderUsecase(repos *OrderUsecaseRepos, configs *OrderUsecaseConfig, logger *zap.Logger) *OrderUsecase {
	return &OrderUsecase{
		repos:   repos,
		configs: configs,
		logger:  logger,
	}
}

func (o *OrderUsecase) CreateOrder(ctx context.Context, params *entity.CreateOrderRequest) error {
	// Validation struct
	if err := libvalidate.Validator().Struct(params); err != nil {
		return libvalidate.ResolveError(err, entity.ErrorCodeInvalidBodyJSON)
	}

	// Retrieve products
	productIDsMap := make(map[string]struct{})
	productIDs := []string{}

	for _, op := range params.Products {
		if _, exists := productIDsMap[op.ProductID]; !exists {
			productIDsMap[op.ProductID] = struct{}{}
			productIDs = append(productIDs, op.ProductID)
		}
	}

	products, err := o.repos.ProductRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return liberr.ResolveError(err)
	}

	// Validate existance products
	if len(productIDs) != len(products) {
		return liberr.ResolveError(entity.ErrorProductNotFound)
	}

	// map[product_id]product
	productMap := map[string]*entity.Product{}
	for _, p := range products {
		productMap[p.ID] = p
	}

	// Retrieve warehouse stocks
	warehouseStocks, err := o.repos.WarehouseRepo.ActiveStock(ctx, productIDs)
	if err != nil {
		return liberr.ResolveError(err)
	}

	// map[product_id][warehouse_id]warehouse_stock
	productWarehouseStockMap := map[string]map[string]*entity.WarehouseStock{}
	for _, ws := range warehouseStocks {
		if _, ok := productWarehouseStockMap[ws.ProductID]; !ok {
			productWarehouseStockMap[ws.ProductID] = map[string]*entity.WarehouseStock{}
		}
		productWarehouseStockMap[ws.ProductID][ws.WarehouseID] = ws
	}

	// Validate warehouse stock
	for _, op := range params.Products {
		if warehouseStock, ok := productWarehouseStockMap[op.ProductID][op.WarehouseID]; !ok {
			return liberr.ResolveError(entity.ErrorProductStockNotFound)
		} else if warehouseStock.ShopID != params.ShopID {
			return liberr.ResolveError(entity.ErrorProductMultiShop)
		} else if warehouseStock.Stock < op.Stock {
			return liberr.ResolveError(entity.ErrorProductInsufficientStock)
		}
	}

	totalStock := 0
	totalPrice := decimal.NewFromInt(0)

	for _, op := range params.Products {
		totalStock += op.Stock

		if product, ok := productMap[op.ProductID]; ok {
			totalPrice = totalPrice.Add(product.Price.Mul(decimal.NewFromInt(int64(op.Stock))))
		}
	}

	tx, err := o.repos.DatabaseTransactionHandler.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return liberr.ResolveError(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback() //nolint
		}
	}()

	now := util.NowUTCWithoutNanoSecond()
	order := &entity.Order{
		UserID:     params.User.ID,
		ShopID:     params.ShopID,
		TotalStock: totalStock,
		TotalPrice: totalPrice,
		ExpiredAt:  now.Add(time.Duration(o.configs.OrderExpirationTimeSecond) * time.Second),
	}

	err = o.repos.OrderRepo.Create(ctx, order, tx)
	if err != nil {
		return liberr.ResolveError(err)
	}

	orderDetail := []*entity.OrderDetail{}
	for _, op := range params.Products {
		price := decimal.NewFromInt(0)
		if product, ok := productMap[op.ProductID]; ok {
			price = product.Price
		}

		orderDetail = append(orderDetail, &entity.OrderDetail{
			OrderID:     order.ID,
			ProductID:   op.ProductID,
			WarehouseID: op.WarehouseID,
			Stock:       op.Stock,
			Price:       price,
		})
	}

	for _, od := range orderDetail {
		err = o.repos.OrderDetailRepo.Create(ctx, od, tx)
		if err != nil {
			return liberr.ResolveError(err)
		}
	}

	err = o.reserveStocks(ctx, params.Products)
	if err != nil {
		return liberr.ResolveError(err)
	}

	err = tx.Commit()
	if err != nil {
		return liberr.ResolveError(err)
	}

	return nil
}

func (o *OrderUsecase) reserveStocks(ctx context.Context, orderProducts []*entity.CreateOrderProduct) error {
	adjustmentStock := []*entity.WarehouseStockAdjustment{}

	for _, p := range orderProducts {
		adjustmentStock = append(adjustmentStock, &entity.WarehouseStockAdjustment{
			WarehouseID: p.WarehouseID,
			ProductID:   p.ProductID,
			Stock:       -1 * p.Stock,
		})
	}

	return o.repos.WarehouseRepo.AdjustmentStock(ctx, &entity.WarehouseStockAdjustmentParams{
		WarehouseStocks: adjustmentStock,
	})
}

func (o *OrderUsecase) releaseStocks(ctx context.Context, orderDetails []*entity.OrderDetail) error {
	adjustmentStock := []*entity.WarehouseStockAdjustment{}

	for _, p := range orderDetails {
		adjustmentStock = append(adjustmentStock, &entity.WarehouseStockAdjustment{
			WarehouseID: p.WarehouseID,
			ProductID:   p.ProductID,
			Stock:       p.Stock,
		})
	}

	return o.repos.WarehouseRepo.AdjustmentStock(ctx, &entity.WarehouseStockAdjustmentParams{
		WarehouseStocks: adjustmentStock,
	})
}
