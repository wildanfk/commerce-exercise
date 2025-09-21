package repository

import (
	"context"
	"fmt"
	"order-service/internal/util"
	"order-service/internal/util/liberr"
	"order-service/module/order/entity"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

var (
	orderDetailTable = "order_details"

	orderDetailInsertColumns = []string{"order_id", "product_id", "warehouse_id", "stock", "price"}
	orderDetailColumns       = []string{"id", "order_id", "product_id", "warehouse_id", "stock", "price", "created_at", "updated_at"}
)

type OrderDetailRepository struct {
	db *sqlx.DB
}

type orderDetailObject struct {
	ID          string          `db:"id"`
	OrderID     string          `db:"order_id"`
	ProductID   string          `db:"product_id"`
	WarehouseID string          `db:"warehouse_id"`
	Stock       int             `db:"stock"`
	Price       decimal.Decimal `db:"price"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

func (od *orderDetailObject) toEntity() *entity.OrderDetail {
	return &entity.OrderDetail{
		ID:          od.ID,
		OrderID:     od.OrderID,
		ProductID:   od.ProductID,
		WarehouseID: od.WarehouseID,
		Stock:       od.Stock,
		Price:       od.Price,
		CreatedAt:   od.CreatedAt,
		UpdatedAt:   od.UpdatedAt,
	}
}

func NewOrderDetailRepository(db *sqlx.DB) *OrderDetailRepository {
	return &OrderDetailRepository{db: db}
}

func (o *OrderDetailRepository) Create(ctx context.Context, orderDetail *entity.OrderDetail, tx util.DatabaseTransaction) error {
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(orderDetailTable)
	ib.Cols(orderDetailInsertColumns...)
	ib.Values(
		orderDetail.OrderID,
		orderDetail.ProductID,
		orderDetail.WarehouseID,
		orderDetail.Stock,
		orderDetail.Price,
	)

	query, args := ib.Build()

	db, err := util.GetExecer(o.db, tx)
	if err != nil {
		return liberr.NewTracer("Error when GetExecer on orderDetail.Create").Wrap(err)
	}

	row, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return liberr.NewTracer("Error when ExecContext on orderDetail.Create").Wrap(err)
	}

	lastInsertedID, err := row.LastInsertId()
	if err != nil {
		return liberr.NewTracer("Error when retrieve LastInsertId on orderDetail.Create").Wrap(err)
	}

	orderDetail.ID = fmt.Sprintf("%d", lastInsertedID)
	return nil
}

func (od *OrderDetailRepository) ListByOrderID(ctx context.Context, orderID string) ([]*entity.OrderDetail, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(orderDetailColumns...)
	sb.From(orderDetailTable)
	sb.Where(sb.Equal("order_id", orderID))

	query, args := sb.Build()

	rows, err := od.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, liberr.NewTracer("Error when QueryxContext on orderDetail.ListByOrderID").Wrap(err)
	}

	orderDetails := []*entity.OrderDetail{}
	for rows.Next() {
		var obj orderDetailObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, liberr.NewTracer("Error when StructScan on orderDetail.ListByOrderID").Wrap(err)
		}

		orderDetails = append(orderDetails, obj.toEntity())
	}

	return orderDetails, nil
}
