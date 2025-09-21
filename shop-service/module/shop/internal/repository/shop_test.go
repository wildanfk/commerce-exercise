package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"shop-service/internal/testutil"
	"shop-service/internal/util/libpagination"
	"shop-service/module/shop/entity"
	"shop-service/module/shop/internal/repository"
	"shop-service/module/shop/testutil/fixtures"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	shopAllAttributes = []string{
		"id",
		"name",
		"created_at",
		"updated_at",
	}

	shopAllColumnsStr = strings.Join(shopAllAttributes, ", ")
)

func TestShopRepository_ListByParams(t *testing.T) {
	columns := shopAllColumnsStr
	rows := shopAllAttributes
	dummyShop := fixtures.NewShop(fixtures.Shop)

	type input struct {
		ctx    context.Context
		params *entity.ListShopByParams
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.Shop, *libpagination.OffsetPagination, error)
	}{
		{
			name: "Success on Retrieve List By Params",
			in: input{
				ctx: context.TODO(),
				params: &entity.ListShopByParams{
					Offset: 5,
					Limit:  20,
					IDs:    []string{"1", "2"},
				},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM shops WHERE id IN (?, ?) LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.IDs[0], in.params.IDs[1], in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetShopRow(dummyShop)...),
					).RowsWillBeClosed()

				expectedCountQuery := "SELECT COUNT(id) AS total FROM shops WHERE id IN (?, ?)"
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
					WithArgs(in.params.IDs[0], in.params.IDs[1]).
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(100)).
					RowsWillBeClosed()

			},
			assertFn: func(result []*entity.Shop, pagination *libpagination.OffsetPagination, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Shop{dummyShop}, result)
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
				params: &entity.ListShopByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM shops LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetShopRow(dummyShop)...),
					).RowsWillBeClosed()

				expectedCountQuery := "SELECT COUNT(id) AS total FROM shops"
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(100)).
					RowsWillBeClosed()

			},
			assertFn: func(result []*entity.Shop, pagination *libpagination.OffsetPagination, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Shop{dummyShop}, result)
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
				params: &entity.ListShopByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM shops LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetShopRow(dummyShop)...),
					).RowsWillBeClosed()

				expectedCountQuery := "SELECT COUNT(id) AS total FROM shops"
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedCountQuery)).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Shop, pagination *libpagination.OffsetPagination, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Nil(t, pagination)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx:    context.TODO(),
				params: &entity.ListShopByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM shops LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(dummyShop.ID, dummyShop.Name, dummyShop.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Shop, pagination *libpagination.OffsetPagination, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Nil(t, pagination)
			},
		},
		{
			name: "Error on QueryxContext",
			in: input{
				ctx:    context.TODO(),
				params: &entity.ListShopByParams{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM shops LIMIT ? OFFSET ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.params.Limit, in.params.Offset).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Shop, pagination *libpagination.OffsetPagination, err error) {
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
			repo := repository.NewShopRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListByParams(tc.in.ctx, tc.in.params))
		})
	}
}
