package repository

import (
	"context"
	"shop-service/internal/util/liberr"
	"shop-service/internal/util/libpagination"
	"shop-service/module/shop/entity"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

var (
	shopTable = "shops"

	shopColumns = []string{"id", "name", "created_at", "updated_at"}
)

type ShopRepository struct {
	db *sqlx.DB
}

type shopObject struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (o *shopObject) toEntity() *entity.Shop {
	return &entity.Shop{
		ID:        o.ID,
		Name:      o.Name,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func NewShopRepository(db *sqlx.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

func (s *ShopRepository) filterByParams(sb *sqlbuilder.SelectBuilder, params *entity.ListShopByParams) *sqlbuilder.SelectBuilder {
	if len(params.IDs) > 0 {
		inArgs := make([]any, len(params.IDs))
		for i, v := range params.IDs {
			inArgs[i] = v
		}
		sb.Where(sb.In("id", inArgs...))
	}

	return sb
}

func (s *ShopRepository) ListByParams(ctx context.Context, params *entity.ListShopByParams) ([]*entity.Shop, *libpagination.OffsetPagination, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(shopColumns...)
	sb.From(shopTable)
	sb.Limit(params.Limit)
	sb.Offset(params.Offset)

	query, args := s.filterByParams(sb, params).Build()

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, nil, liberr.NewTracer("Error when QueryxContext on user.ListByParams").Wrap(err)
	}

	shops := []*entity.Shop{}
	for rows.Next() {
		var obj shopObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, nil, liberr.NewTracer("Error when StructScan on user.ListByParams").Wrap(err)
		}

		shops = append(shops, obj.toEntity())
	}

	cb := sqlbuilder.NewSelectBuilder()
	cb.Select(cb.As("COUNT(id)", "total"))
	cb.From(shopTable)

	cQuery, cArgs := s.filterByParams(cb, params).Build()
	row := s.db.QueryRowxContext(ctx, cQuery, cArgs...)

	var total int
	if err := row.Scan(&total); err != nil {
		return nil, nil, liberr.NewTracer("Error when Scan on user.ListByParams: %s").Wrap(err)
	}

	return shops, &libpagination.OffsetPagination{
		Total:  total,
		Offset: params.Offset,
		Limit:  params.Limit,
	}, nil
}
