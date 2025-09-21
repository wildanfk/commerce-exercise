package repository

import (
	"context"
	"time"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/module/warehouse/entity"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

var (
	warehouseTable = "warehouses"

	warehouseColumns = []string{"id", "shop_id", "name", "active", "created_at", "updated_at"}
)

type WarehouseRepository struct {
	db *sqlx.DB
}

type warehouseObject struct {
	ID        string    `db:"id"`
	ShopID    string    `db:"shop_id"`
	Name      string    `db:"name"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (o *warehouseObject) toEntity() *entity.Warehouse {
	return &entity.Warehouse{
		ID:        o.ID,
		ShopID:    o.ShopID,
		Name:      o.Name,
		Active:    o.Active,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func NewWarehouseRepository(db *sqlx.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (w *WarehouseRepository) ListByIDs(ctx context.Context, ids []string) ([]*entity.Warehouse, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(warehouseColumns...)
	sb.From(warehouseTable)

	if len(ids) > 0 {
		inArgs := make([]any, len(ids))
		for i, v := range ids {
			inArgs[i] = v
		}
		sb.Where(sb.In("id", inArgs...))
	}

	query, args := sb.Build()

	rows, err := w.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, liberr.NewTracer("Error when QueryxContext on user.ListByIDs").Wrap(err)
	}

	warehouses := []*entity.Warehouse{}
	for rows.Next() {
		var obj warehouseObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, liberr.NewTracer("Error when StructScan on user.ListByIDs").Wrap(err)
		}

		warehouses = append(warehouses, obj.toEntity())
	}

	return warehouses, nil
}

func (w *WarehouseRepository) UpdateActive(ctx context.Context, id string, active bool) error {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(warehouseTable).
		Set(
			ub.Assign("active", active),
		).
		Where(
			ub.E("id", id),
		)
	query, args := ub.Build()

	_, err := w.db.ExecContext(ctx, query, args...)
	if err != nil {
		return liberr.NewTracer("Error when ExecContext on user.UpdateActive").Wrap(err)
	}

	return nil
}
