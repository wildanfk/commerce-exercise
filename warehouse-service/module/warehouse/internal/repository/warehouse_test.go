package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"warehouse-service/internal/testutil"
	"warehouse-service/module/warehouse/entity"
	"warehouse-service/module/warehouse/internal/repository"
	"warehouse-service/module/warehouse/testutil/fixtures"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	warehouseAllAttributes = []string{
		"id",
		"shop_id",
		"name",
		"active",
		"created_at",
		"updated_at",
	}

	warehouseAllColumnsStr = strings.Join(warehouseAllAttributes, ", ")
)

func TestWarehouseRepository_ListByIDs(t *testing.T) {
	columns := warehouseAllColumnsStr
	rows := warehouseAllAttributes
	dummyWarehouse := fixtures.NewWarehouse(fixtures.Warehouse)

	type input struct {
		ctx context.Context
		ids []string
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.Warehouse, error)
	}{
		{
			name: "Success on Retrieve ListByIDs",
			in: input{
				ctx: context.TODO(),
				ids: []string{"1", "2"},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouses WHERE id IN (?, ?)", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs("1", "2").
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetWarehouseRow(dummyWarehouse)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Warehouse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Warehouse{dummyWarehouse}, result)
			},
		},
		{
			name: "Success on Retrieve ListByIDs With Empty Params",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouses", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetWarehouseRow(dummyWarehouse)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Warehouse, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Warehouse{dummyWarehouse}, result)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouses", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(dummyWarehouse.ID, dummyWarehouse.ShopID, dummyWarehouse.Name, dummyWarehouse.Active, dummyWarehouse.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Warehouse, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "Error on QueryxContext",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouses", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Warehouse, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListByIDs(tc.in.ctx, tc.in.ids))
		})
	}
}

func TestWarehouseRepository_UpdateActive(t *testing.T) {
	expectedQuery := "UPDATE warehouses SET active = ? WHERE id  = ?"

	type input struct {
		ctx    context.Context
		id     string
		active bool
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func(error)
	}{
		{
			name: "Success on Update",
			in: input{
				ctx:    context.TODO(),
				id:     "1",
				active: true,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(true, "1").
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(nil)
			},
			assertFn: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx:    context.TODO(),
				id:     "1",
				active: true,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(true, "1").
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(errors.New("error"))
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.UpdateActive(tc.in.ctx, tc.in.id, tc.in.active))
		})
	}
}
