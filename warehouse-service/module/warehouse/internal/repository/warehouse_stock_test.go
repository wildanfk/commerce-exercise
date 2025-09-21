package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"warehouse-service/internal/util"
	"warehouse-service/module/warehouse/entity"
	"warehouse-service/module/warehouse/internal/repository"
	"warehouse-service/module/warehouse/testutil/fixtures"

	"warehouse-service/internal/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	warehouseStockInsertAttributes = []string{
		"warehouse_id",
		"product_id",
		"stock",
	}
	warehouseStockAllAttributes = []string{
		"id",
		"warehouse_id",
		"product_id",
		"stock",
		"created_at",
		"updated_at",
	}

	warehouseStockInsertColumnsStr = strings.Join(warehouseStockInsertAttributes, ", ")
	warehouseStockAllColumnsStr    = strings.Join(warehouseStockAllAttributes, ", ")
)

func TestWarehouseStockRepository_Create(t *testing.T) {
	expectedQuery := fmt.Sprintf("INSERT INTO warehouse_stocks (%s) VALUES (?, ?, ?)", warehouseStockInsertColumnsStr)

	type input struct {
		ctx            context.Context
		warehouseStock *entity.WarehouseStock
		tx             util.DatabaseTransaction
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func(error)
	}{
		{
			name: "Success on Create",
			in: input{
				ctx:            context.TODO(),
				warehouseStock: fixtures.NewWarehouseStock(fixtures.WarehouseStock),
				tx:             nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.warehouseStock.WarehouseID, in.warehouseStock.ProductID, in.warehouseStock.Stock).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(nil)
			},
			assertFn: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "Error on LastInsertId",
			in: input{
				ctx:            context.TODO(),
				warehouseStock: fixtures.NewWarehouseStock(fixtures.WarehouseStock),
				tx:             nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.warehouseStock.WarehouseID, in.warehouseStock.ProductID, in.warehouseStock.Stock).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx:            context.TODO(),
				warehouseStock: fixtures.NewWarehouseStock(fixtures.WarehouseStock),
				tx:             nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.warehouseStock.WarehouseID, in.warehouseStock.ProductID, in.warehouseStock.Stock).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(errors.New("error"))
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error on Execute Query Duplicate Error Happened",
			in: input{
				ctx:            context.TODO(),
				warehouseStock: fixtures.NewWarehouseStock(fixtures.WarehouseStock),
				tx:             nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.warehouseStock.WarehouseID, in.warehouseStock.ProductID, in.warehouseStock.Stock).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error on GetExecer",
			in: input{
				ctx:            context.TODO(),
				warehouseStock: fixtures.NewWarehouseStock(fixtures.WarehouseStock),
				tx:             &testutil.UnknownDatabaseTransaction{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseStockRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.Create(tc.in.ctx, tc.in.warehouseStock, tc.in.tx))
		})
	}
}

func TestWarehouseStockRepository_IncreaseStock(t *testing.T) {
	expectedQuery := "UPDATE warehouse_stocks SET stock = stock + ? WHERE warehouse_id  = ? AND product_id = ?"

	type input struct {
		ctx    context.Context
		params entity.WarehouseStockAdjustmentParams
		tx     util.DatabaseTransaction
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func(int64, error)
	}{
		{
			name: "Success on Update",
			in: input{
				ctx: context.TODO(),
				params: entity.WarehouseStockAdjustmentParams{
					WarehouseID: "1",
					ProductID:   "2",
					Stock:       10,
				},
				tx: nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(10, "1", "2").
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(nil)
			},
			assertFn: func(result int64, err error) {
				assert.Nil(t, err)
				assert.Equal(t, int64(1), result)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx: context.TODO(),
				params: entity.WarehouseStockAdjustmentParams{
					WarehouseID: "1",
					ProductID:   "2",
					Stock:       10,
				},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(10, "1", "2").
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(errors.New("error"))
			},
			assertFn: func(result int64, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, int64(0), result)
			},
		},
		{
			name: "Error on GetExecer",
			in: input{
				ctx: context.TODO(),
				params: entity.WarehouseStockAdjustmentParams{
					WarehouseID: "1",
					ProductID:   "2",
					Stock:       10,
				},
				tx: &testutil.UnknownDatabaseTransaction{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {},
			assertFn: func(result int64, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, int64(0), result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseStockRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.IncreaseStock(tc.in.ctx, tc.in.params, tc.in.tx))
		})
	}
}

func TestWarehouseStockRepository_DecreaseStock(t *testing.T) {
	expectedQuery := "UPDATE warehouse_stocks SET stock = stock - ? WHERE warehouse_id  = ? AND product_id = ? AND stock >= ?"

	type input struct {
		ctx    context.Context
		params entity.WarehouseStockAdjustmentParams
		tx     util.DatabaseTransaction
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func(int64, error)
	}{
		{
			name: "Success on Update",
			in: input{
				ctx: context.TODO(),
				params: entity.WarehouseStockAdjustmentParams{
					WarehouseID: "1",
					ProductID:   "2",
					Stock:       10,
				},
				tx: nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(10, "1", "2", 10).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(nil)
			},
			assertFn: func(result int64, err error) {
				assert.Nil(t, err)
				assert.Equal(t, int64(1), result)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx: context.TODO(),
				params: entity.WarehouseStockAdjustmentParams{
					WarehouseID: "1",
					ProductID:   "2",
					Stock:       10,
				},
				tx: nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(10, "1", "2", 10).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(errors.New("error"))
			},
			assertFn: func(result int64, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, int64(0), result)
			},
		},
		{
			name: "Error on GetExecer",
			in: input{
				ctx: context.TODO(),
				params: entity.WarehouseStockAdjustmentParams{
					WarehouseID: "1",
					ProductID:   "2",
					Stock:       10,
				},
				tx: &testutil.UnknownDatabaseTransaction{},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {},
			assertFn: func(result int64, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, int64(0), result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseStockRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.DecreaseStock(tc.in.ctx, tc.in.params, tc.in.tx))
		})
	}
}

func TestWarehouseStockRepository_ListByWarehouseIDsAndProductIDs(t *testing.T) {
	columns := warehouseStockAllColumnsStr
	rows := warehouseStockAllAttributes
	dummyWarehouseStock := fixtures.NewWarehouseStock(fixtures.WarehouseStock)

	type input struct {
		ctx          context.Context
		warehouseIDs []string
		productIDs   []string
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.WarehouseStock, error)
	}{
		{
			name: "Success on Retrieve ListByWarehouseIDAndProductIDs",
			in: input{
				ctx:          context.TODO(),
				warehouseIDs: []string{"10"},
				productIDs:   []string{"1", "2"},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks WHERE warehouse_id IN (?) AND product_id IN (?, ?)", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs("10", "1", "2").
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetWarehouseStockRow(dummyWarehouseStock)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.WarehouseStock{dummyWarehouseStock}, result)
			},
		},
		{
			name: "Success on Retrieve ListByWarehouseIDAndProductIDs With Empty Params",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetWarehouseStockRow(dummyWarehouseStock)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.WarehouseStock{dummyWarehouseStock}, result)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(
								dummyWarehouseStock.ID,
								dummyWarehouseStock.WarehouseID, dummyWarehouseStock.ProductID,
								dummyWarehouseStock.Stock, dummyWarehouseStock.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
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
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseStockRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListByWarehouseIDsAndProductIDs(tc.in.ctx, tc.in.warehouseIDs, tc.in.productIDs))
		})
	}
}

func TestWarehouseStockRepository_ListActiveByProductIDs(t *testing.T) {
	columns := warehouseStockAllColumnsStr
	rows := warehouseStockAllAttributes
	dummyWarehouseStock := fixtures.NewWarehouseStock(fixtures.WarehouseStock)

	type input struct {
		ctx        context.Context
		productIDs []string
	}

	activeWarehouseQuery := "SELECT id FROM warehouses WHERE active = ?"

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.WarehouseStock, error)
	}{
		{
			name: "Success on Retrieve ListActiveByProductIDs",
			in: input{
				ctx:        context.TODO(),
				productIDs: []string{"1", "2"},
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks WHERE product_id IN (?, ?) AND warehouse_id IN (%s)", columns, activeWarehouseQuery)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs("1", "2", true).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetWarehouseStockRow(dummyWarehouseStock)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.WarehouseStock{dummyWarehouseStock}, result)
			},
		},
		{
			name: "Success on Retrieve ListActiveByProductIDs With Empty Params",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks WHERE warehouse_id IN (%s)", columns, activeWarehouseQuery)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(true).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetWarehouseStockRow(dummyWarehouseStock)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.WarehouseStock{dummyWarehouseStock}, result)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks WHERE warehouse_id IN (%s)", columns, activeWarehouseQuery)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(true).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(
								dummyWarehouseStock.ID,
								dummyWarehouseStock.WarehouseID, dummyWarehouseStock.ProductID,
								dummyWarehouseStock.Stock, dummyWarehouseStock.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
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
				expectedQuery := fmt.Sprintf("SELECT %s FROM warehouse_stocks WHERE warehouse_id IN (%s)", columns, activeWarehouseQuery)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(true).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.WarehouseStock, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewWarehouseStockRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListActiveByProductIDs(tc.in.ctx, tc.in.productIDs))
		})
	}
}
