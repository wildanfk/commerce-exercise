package repository_test

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/testutil"
	"order-service/internal/util"
	"order-service/module/order/entity"
	"order-service/module/order/internal/repository"
	"order-service/module/order/testutil/fixtures"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	orderInsertAttributes = []string{
		"user_id",
		"shop_id",
		"state",
		"total_stock",
		"total_price",
		"expired_at",
	}
	orderAllAttributes = []string{
		"id",
		"user_id",
		"shop_id",
		"state",
		"total_stock",
		"total_price",
		"expired_at",
		"created_at",
		"updated_at",
	}

	orderInsertColumnsStr = strings.Join(orderInsertAttributes, ", ")
	orderAllColumnsStr    = strings.Join(orderAllAttributes, ", ")
)

func TestOrderRepository_Create(t *testing.T) {
	expectedQuery := fmt.Sprintf("INSERT INTO orders (%s) VALUES (?, ?, ?, ?, ?, ?)", orderInsertColumnsStr)

	type input struct {
		ctx   context.Context
		order *entity.Order
		tx    util.DatabaseTransaction
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
				ctx:   context.TODO(),
				order: fixtures.NewOrder(fixtures.Order),
				tx:    nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.order.UserID, in.order.ShopID, entity.OrderStateCreated, in.order.TotalStock, in.order.TotalPrice, in.order.ExpiredAt).
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
				ctx:   context.TODO(),
				order: fixtures.NewOrder(fixtures.Order),
				tx:    nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.order.UserID, in.order.ShopID, entity.OrderStateCreated, in.order.TotalStock, in.order.TotalPrice, in.order.ExpiredAt).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx:   context.TODO(),
				order: fixtures.NewOrder(fixtures.Order),
				tx:    nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.order.UserID, in.order.ShopID, entity.OrderStateCreated, in.order.TotalStock, in.order.TotalPrice, in.order.ExpiredAt).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(errors.New("error"))
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error on GetExecer",
			in: input{
				ctx:   context.TODO(),
				order: fixtures.NewOrder(fixtures.Order),
				tx:    &testutil.UnknownDatabaseTransaction{},
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
			repo := repository.NewOrderRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.Create(tc.in.ctx, tc.in.order, tc.in.tx))
		})
	}
}

func TestOrderRepository_UpdateExpired(t *testing.T) {
	expectedQuery := "UPDATE orders SET state = ? WHERE id = ?"

	type input struct {
		ctx context.Context
		id  string
		tx  util.DatabaseTransaction
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
				id:  "1",
				tx:  nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(entity.OrderStateExpired, "1").
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
				id:  "1",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(entity.OrderStateExpired, "1").
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
				id:  "1",
				tx:  &testutil.UnknownDatabaseTransaction{},
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
			repo := repository.NewOrderRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.UpdateExpired(tc.in.ctx, tc.in.id, tc.in.tx))
		})
	}
}

func TestOrderRepository_ListByOrderExpired(t *testing.T) {
	columns := orderAllColumnsStr
	rows := orderAllAttributes
	dummyOrder := fixtures.NewOrder(fixtures.Order)

	type input struct {
		ctx context.Context
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.Order, error)
	}{
		{
			name: "Success on Retrieve ListByOrderExpired",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM orders WHERE state = ? AND expired_at < NOW()", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(entity.OrderStateCreated).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetOrderRow(dummyOrder)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Order, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.Order{dummyOrder}, result)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx: context.TODO(),
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM orders WHERE state = ? AND expired_at < NOW()", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(entity.OrderStateCreated).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(
								dummyOrder.ID,
								dummyOrder.UserID, dummyOrder.ShopID, dummyOrder.State,
								dummyOrder.TotalStock, dummyOrder.TotalPrice,
								dummyOrder.ExpiredAt, dummyOrder.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Order, err error) {
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
				expectedQuery := fmt.Sprintf("SELECT %s FROM orders WHERE state = ? AND expired_at < NOW()", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(entity.OrderStateCreated).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.Order, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewOrderRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListByOrderExpired(tc.in.ctx))
		})
	}
}
