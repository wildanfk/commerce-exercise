package repository

import (
	"context"
	"fmt"
	"time"
	"warehouse-service/internal/util"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/module/warehouse/entity"

	"github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

var (
	warehouseStockTable = "warehouse_stocks"

	warehouseStockInsertColumns = []string{"warehouse_id", "product_id", "stock"}
	warehouseStockColumns       = []string{"id", "warehouse_id", "product_id", "stock", "created_at", "updated_at"}
)

type WarehouseStockRepository struct {
	db *sqlx.DB
}

type warehouseStockObject struct {
	ID          string    `db:"id"`
	WarehouseID string    `db:"warehouse_id"`
	ProductID   string    `db:"product_id"`
	Stock       int       `db:"stock"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (o *warehouseStockObject) toEntity() *entity.WarehouseStock {
	return &entity.WarehouseStock{
		ID:          o.ID,
		WarehouseID: o.WarehouseID,
		ProductID:   o.ProductID,
		Stock:       o.Stock,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func NewWarehouseStockRepository(db *sqlx.DB) *WarehouseStockRepository {
	return &WarehouseStockRepository{db: db}
}

func (w *WarehouseStockRepository) Create(ctx context.Context, warehouseStock *entity.WarehouseStock, tx util.DatabaseTransaction) error {
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(warehouseStockTable)
	ib.Cols(warehouseStockInsertColumns...)
	ib.Values(
		warehouseStock.WarehouseID,
		warehouseStock.ProductID,
		warehouseStock.Stock,
	)
	query, args := ib.Build()

	db, err := util.GetExecer(w.db, tx)
	if err != nil {
		return liberr.NewTracer("Error when GetExecer on user.Create").Wrap(err)
	}

	row, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		if isDuplicateError(err) {
			return entity.ErrorWarehouseStockDuplicated
		}

		return liberr.NewTracer("Error when ExecContext on user.Create").Wrap(err)
	}

	lastInsertedID, err := row.LastInsertId()
	if err != nil {
		return liberr.NewTracer("Error when retrieve LastInsertId on user.Create").Wrap(err)
	}

	warehouseStock.ID = fmt.Sprintf("%d", lastInsertedID)
	return nil
}

func isDuplicateError(err error) bool {
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}

func (w *WarehouseStockRepository) IncreaseStock(ctx context.Context, params entity.WarehouseStockAdjustmentParams, tx util.DatabaseTransaction) (int64, error) {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(warehouseStockTable).
		Set(
			ub.Add("stock", params.Stock),
		).
		Where(
			ub.E("warehouse_id", params.WarehouseID),
			ub.E("product_id", params.ProductID),
		)
	query, args := ub.Build()

	db, err := util.GetExecer(w.db, tx)
	if err != nil {
		return 0, liberr.NewTracer("Error when GetExecer on user.IncreaseStock").Wrap(err)
	}

	row, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, liberr.NewTracer("Error when ExecContext on user.IncreaseStock").Wrap(err)
	}

	rowAffected, _ := row.RowsAffected()
	return rowAffected, nil
}

func (w *WarehouseStockRepository) DecreaseStock(ctx context.Context, params entity.WarehouseStockAdjustmentParams, tx util.DatabaseTransaction) (int64, error) {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(warehouseStockTable).
		Set(
			ub.Sub("stock", params.Stock),
		).
		Where(
			ub.E("warehouse_id", params.WarehouseID),
			ub.E("product_id", params.ProductID),
			ub.GTE("stock", params.Stock),
		)
	query, args := ub.Build()

	db, err := util.GetExecer(w.db, tx)
	if err != nil {
		return 0, liberr.NewTracer("Error when GetExecer on user.DecreaseStock").Wrap(err)
	}

	row, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, liberr.NewTracer("Error when ExecContext on user.DecreaseStock").Wrap(err)
	}

	rowAffected, _ := row.RowsAffected()
	return rowAffected, nil
}

func (w *WarehouseStockRepository) ListByWarehouseIDsAndProductIDs(ctx context.Context, warehouseIDs []string, productIDs []string) ([]*entity.WarehouseStock, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(warehouseStockColumns...)
	sb.From(warehouseStockTable)

	if len(warehouseIDs) > 0 {
		inWarehouseIDsArgs := make([]any, len(warehouseIDs))
		for i, v := range warehouseIDs {
			inWarehouseIDsArgs[i] = v
		}
		sb.Where(sb.In("warehouse_id", inWarehouseIDsArgs...))
	}

	if len(productIDs) > 0 {
		inProductIDsArgs := make([]any, len(productIDs))
		for i, v := range productIDs {
			inProductIDsArgs[i] = v
		}
		sb.Where(sb.In("product_id", inProductIDsArgs...))
	}

	query, args := sb.Build()

	rows, err := w.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, liberr.NewTracer("Error when QueryxContext on user.ListByWarehouseIDAndProductIDs").Wrap(err)
	}

	warehouseStocks := []*entity.WarehouseStock{}
	for rows.Next() {
		var obj warehouseStockObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, liberr.NewTracer("Error when StructScan on user.ListByWarehouseIDAndProductIDs").Wrap(err)
		}

		warehouseStocks = append(warehouseStocks, obj.toEntity())
	}

	return warehouseStocks, nil
}

func (w *WarehouseStockRepository) ListActiveByProductIDs(ctx context.Context, productIDs []string) ([]*entity.WarehouseStock, error) {
	activeWarehouseSb := sqlbuilder.NewSelectBuilder()
	activeWarehouseSb.Select("id")
	activeWarehouseSb.From(warehouseTable)
	activeWarehouseSb.Where(activeWarehouseSb.Equal("active", true))

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(warehouseStockColumns...)
	sb.From(warehouseStockTable)

	if len(productIDs) > 0 {
		inArgs := make([]any, len(productIDs))
		for i, v := range productIDs {
			inArgs[i] = v
		}
		sb.Where(sb.In("product_id", inArgs...))
	}

	sb.Where(sb.In("warehouse_id", activeWarehouseSb))

	query, args := sb.Build()

	rows, err := w.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, liberr.NewTracer("Error when QueryxContext on user.ListActiveByProductIDs").Wrap(err)
	}

	warehouseStocks := []*entity.WarehouseStock{}
	for rows.Next() {
		var obj warehouseStockObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, liberr.NewTracer("Error when StructScan on user.ListActiveByProductIDs").Wrap(err)
		}

		warehouseStocks = append(warehouseStocks, obj.toEntity())
	}

	return warehouseStocks, nil
}
