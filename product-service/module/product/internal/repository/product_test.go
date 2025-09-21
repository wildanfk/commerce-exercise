package repository_test

import (
	"context"
	"fmt"
	"product-service/internal/testutil"
	"product-service/internal/util/libpagination"
	"product-service/module/product/entity"
	"product-service/module/product/internal/repository"
	"product-service/module/product/testutil/fixtures"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	productAllAttributes = []string{
		"id",
		"name",
		"price",
		"created_at",
		"updated_at",
	}

	productAllColumnsStr = strings.Join(productAllAttributes, ", ")
)

func TestShopRepository_ListByParams(t *testing.T) {
	columns := productAllColumnsStr
	rows := productAllAttributes
	dummyProduct := fixtures.NewProduct(fixtures.Product)

	type input struct {
		ctx    context.Context
		params *entity.ListProductByParams
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.Product, *libpagination.OffsetPagination, error)
	}{
		{
			name: "Success on Retrieve List By Params",
			in: input{
				ctx: context.TODO(),
				params: &entity.ListProductByParams{
					Offset: 5,
					Limit:  20,
					IDs:    []string{"1", "2"},
					Name:   "test",
				},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM products WHERE id IN (?, ?) AND name LIKE ? LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.IDs[0], in.params.IDs[1], "test%", in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetProductRow(dummyProduct)...),
					).RowsWillBeClosed()

				expectedCountQuery := "SELECT COUNT(id) AS total FROM products WHERE id IN (?, ?) AND name LIKE ?"
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
					WithArgs(in.params.IDs[0], in.params.IDs[1], "test%").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(100)).
					RowsWillBeClosed()

			},
			assertFn: func(result []*entity.Product, pagination *libpagination.OffsetPagination, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Product{dummyProduct}, result)
				assert.Equal(t, &libpagination.OffsetPagination{
					Offset: 5,
					Limit:  20,
					Total:  100,
				}, pagination)
			},
		},
		{
			name: "Success on Retrieve List By Params With Empty Params",
			in: input{
				ctx:    context.TODO(),
				params: &entity.ListProductByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM products LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetProductRow(dummyProduct)...),
					).RowsWillBeClosed()

				expectedCountQuery := "SELECT COUNT(id) AS total FROM products"
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(100)).
					RowsWillBeClosed()

			},
			assertFn: func(result []*entity.Product, pagination *libpagination.OffsetPagination, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Product{dummyProduct}, result)
				assert.Equal(t, &libpagination.OffsetPagination{
					Offset: 0,
					Limit:  0,
					Total:  100,
				}, pagination)
			},
		},
		{
			name: "Error on Scan Count Query",
			in: input{
				ctx:    context.TODO(),
				params: &entity.ListProductByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM products LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetProductRow(dummyProduct)...),
					).RowsWillBeClosed()

				expectedCountQuery := "SELECT COUNT(id) AS total FROM products"
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Product, pagination *libpagination.OffsetPagination, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Nil(t, pagination)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx:    context.TODO(),
				params: &entity.ListProductByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM products LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(dummyProduct.ID, dummyProduct.Name, dummyProduct.Price, dummyProduct.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Product, pagination *libpagination.OffsetPagination, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Nil(t, pagination)
			},
		},
		{
			name: "Error on QueryxContext",
			in: input{
				ctx:    context.TODO(),
				params: &entity.ListProductByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM products LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Product, pagination *libpagination.OffsetPagination, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Nil(t, pagination)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewProductRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListByParams(tc.in.ctx, tc.in.params))
		})
	}
}
