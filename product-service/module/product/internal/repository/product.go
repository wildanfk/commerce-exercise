package repository

import (
	"context"
	"fmt"
	"product-service/internal/util/liberr"
	"product-service/internal/util/libpagination"
	"product-service/module/product/entity"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

var (
	productTable = "products"

	productColumns = []string{"id", "name", "price", "created_at", "updated_at"}
)

type ProductRepository struct {
	db *sqlx.DB
}

type productObject struct {
	ID        string          `db:"id"`
	Name      string          `db:"name"`
	Price     decimal.Decimal `db:"price"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (o *productObject) toEntity() *entity.Product {
	return &entity.Product{
		ID:        o.ID,
		Name:      o.Name,
		Price:     o.Price,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) filterByParams(sb *sqlbuilder.SelectBuilder, params *entity.ListProductByParams) *sqlbuilder.SelectBuilder {
	if len(params.IDs) > 0 {
		inArgs := make([]any, len(params.IDs))
		for i, v := range params.IDs {
			inArgs[i] = v
		}
		sb.Where(sb.In("id", inArgs...))
	}
	if params.Name != "" {
		sb.Where(sb.Like("name", fmt.Sprintf("%s%%", params.Name)))
	}

	return sb
}

func (s *ProductRepository) ListByParams(ctx context.Context, params *entity.ListProductByParams) ([]*entity.Product, *libpagination.OffsetPagination, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(productColumns...)
	sb.From(productTable)
	sb.Limit(params.Limit)
	sb.Offset(params.Offset)

	query, args := s.filterByParams(sb, params).Build()

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, nil, liberr.NewTracer("Error when QueryxContext on user.ListByParams").Wrap(err)
	}

	products := []*entity.Product{}
	for rows.Next() {
		var obj productObject

		if err := rows.StructScan(&obj); err != nil {
			return nil, nil, liberr.NewTracer("Error when StructScan on user.ListByParams").Wrap(err)
		}

		products = append(products, obj.toEntity())
	}

	cb := sqlbuilder.NewSelectBuilder()
	cb.Select(cb.As("COUNT(id)", "total"))
	cb.From(productTable)

	cQuery, cArgs := s.filterByParams(cb, params).Build()
	row := s.db.QueryRowxContext(ctx, cQuery, cArgs...)

	var total int
	if err := row.Scan(&total); err != nil {
		return nil, nil, liberr.NewTracer("Error when Scan on user.ListByParams: %s").Wrap(err)
	}

	return products, &libpagination.OffsetPagination{
		Total:  total,
		Offset: params.Offset,
		Limit:  params.Limit,
	}, nil
}
