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
	orderTable = "orders"

	orderInsertColumns = []string{"user_id", "shop_id", "state", "total_stock", "total_price", "expired_at"}
	orderColumns       = []string{"id", "user_id", "shop_id", "state", "total_stock", "total_price", "expired_at", "created_at", "updated_at"}
)

type OrderRepository struct {
	db *sqlx.DB
}

type orderObject struct {
	ID         string          `db:"id"`
	UserID     string          `db:"user_id"`
	ShopID     string          `db:"shop_id"`
	State      int             `db:"state"`
	TotalStock int             `db:"total_stock"`
	TotalPrice decimal.Decimal `db:"total_price"`
	ExpiredAt  time.Time       `db:"expired_at"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}

func (o *orderObject) toEntity() *entity.Order {
	return &entity.Order{
		ID:         o.ID,
		UserID:     o.UserID,
		ShopID:     o.ShopID,
		State:      entity.OrderState(o.State),
		TotalStock: o.TotalStock,
		TotalPrice: o.TotalPrice,
		ExpiredAt:  o.ExpiredAt,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) Create(ctx context.Context, order *entity.Order, tx util.DatabaseTransaction) error {
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(orderTable)
	ib.Cols(orderInsertColumns...)
	ib.Values(
		order.UserID,
		order.ShopID,
		entity.OrderStateCreated,
		order.TotalStock,
		order.TotalPrice,
		order.ExpiredAt,
	)

	query, args := ib.Build()

	db, err := util.GetExecer(o.db, tx)
	if err != nil {
		return liberr.NewTracer("Error when GetExecer on order.Create").Wrap(err)
	}

	row, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return liberr.NewTracer("Error when ExecContext on order.Create").Wrap(err)
	}

	lastInsertedID, err := row.LastInsertId()
	if err != nil {
		return liberr.NewTracer("Error when retrieve LastInsertId on order.Create").Wrap(err)
	}

	order.ID = fmt.Sprintf("%d", lastInsertedID)
	return nil
}

func (o *OrderRepository) UpdateExpired(ctx context.Context, id string, tx util.DatabaseTransaction) (int64, error) {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(orderTable).
		Set(
			ub.Assign("state", entity.OrderStateExpired),
		).
		Where(
			ub.E("id", id),
		)
	query, args := ub.Build()

	db, err := util.GetExecer(o.db, tx)
	if err != nil {
		return 0, liberr.NewTracer("Error when GetExecer on order.UpdateExpired").Wrap(err)
	}

	row, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, liberr.NewTracer("Error when ExecContext on order.UpdateExpired").Wrap(err)
	}

	rowAffected, _ := row.RowsAffected()
	return rowAffected, nil
}

func (o *OrderRepository) ListByOrderExpired(ctx context.Context) ([]*entity.Order, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(orderColumns...)
	sb.From(orderTable)
	sb.Where(sb.Equal("state", entity.OrderStateCreated))
	sb.Where("expired_at < NOW()")

	query, args := sb.Build()

	rows, err := o.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, liberr.NewTracer("Error when QueryxContext on order.ListByOrderExpired").Wrap(err)
	}

	orders := []*entity.Order{}
	for rows.Next() {
		var obj orderObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, liberr.NewTracer("Error when StructScan on user.ListByWarehouseIDAndProductIDs").Wrap(err)
		}

		orders = append(orders, obj.toEntity())
	}

	return orders, nil
}
