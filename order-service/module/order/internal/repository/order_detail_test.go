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
	orderDetailInsertAttributes = []string{
		"order_id",
		"product_id",
		"warehouse_id",
		"stock",
		"price",
	}
	orderDetailAllAttributes = []string{
		"id",
		"order_id",
		"product_id",
		"warehouse_id",
		"stock",
		"price",
		"created_at",
		"updated_at",
	}

	orderDetailInsertColumnsStr = strings.Join(orderDetailInsertAttributes, ", ")
	orderDetailAllColumnsStr    = strings.Join(orderDetailAllAttributes, ", ")
)

func TestOrderDetailRepository_Create(t *testing.T) {
	expectedQuery := fmt.Sprintf("INSERT INTO order_details (%s) VALUES (?, ?, ?, ?, ?)", orderDetailInsertColumnsStr)

	type input struct {
		ctx         context.Context
		orderDetail *entity.OrderDetail
		tx          util.DatabaseTransaction
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
				ctx:         context.TODO(),
				orderDetail: fixtures.NewOrderDetail(fixtures.OrderDetail),
				tx:          nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.orderDetail.OrderID, in.orderDetail.ProductID, in.orderDetail.WarehouseID, in.orderDetail.Stock, in.orderDetail.Price).
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
				ctx:         context.TODO(),
				orderDetail: fixtures.NewOrderDetail(fixtures.OrderDetail),
				tx:          nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.orderDetail.OrderID, in.orderDetail.ProductID, in.orderDetail.WarehouseID, in.orderDetail.Stock, in.orderDetail.Price).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx:         context.TODO(),
				orderDetail: fixtures.NewOrderDetail(fixtures.OrderDetail),
				tx:          nil,
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := regexp.QuoteMeta(expectedQuery)
				dependency.MockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(in.orderDetail.OrderID, in.orderDetail.ProductID, in.orderDetail.WarehouseID, in.orderDetail.Stock, in.orderDetail.Price).
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
				ctx:         context.TODO(),
				orderDetail: fixtures.NewOrderDetail(fixtures.OrderDetail),
				tx:          &testutil.UnknownDatabaseTransaction{},
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
			repo := repository.NewOrderDetailRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.Create(tc.in.ctx, tc.in.orderDetail, tc.in.tx))
		})
	}
}

func TestOrderDetailRepository_ListByOrderID(t *testing.T) {
	columns := orderDetailAllColumnsStr
	rows := orderDetailAllAttributes
	dummyOrderDetail := fixtures.NewOrderDetail(fixtures.OrderDetail)

	type input struct {
		ctx     context.Context
		orderID string
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func([]*entity.OrderDetail, error)
	}{
		{
			name: "Success on Retrieve ListByOrderExpired",
			in: input{
				ctx:     context.TODO(),
				orderID: "1",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM order_details WHERE order_id = ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs("1").
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetOrderDetailRow(dummyOrderDetail)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.OrderDetail, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []*entity.OrderDetail{dummyOrderDetail}, result)
			},
		},
		{
			name: "Error on StructScan",
			in: input{
				ctx:     context.TODO(),
				orderID: "1",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM order_details WHERE order_id = ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs("1").
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(
								dummyOrderDetail.ID,
								dummyOrderDetail.OrderID, dummyOrderDetail.ProductID, dummyOrderDetail.WarehouseID,
								dummyOrderDetail.Stock, dummyOrderDetail.Price,
								dummyOrderDetail.CreatedAt, "invalid"),
					).RowsWillBeClosed()
			},
			assertFn: func(result []*entity.OrderDetail, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "Error on QueryxContext",
			in: input{
				ctx:     context.TODO(),
				orderID: "1",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				expectedQuery := fmt.Sprintf("SELECT %s FROM order_details WHERE order_id = ?", columns)
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs("1").
					WithArgs(entity.OrderStateCreated).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result []*entity.OrderDetail, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewOrderDetailRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.ListByOrderID(tc.in.ctx, tc.in.orderID))
		})
	}
}
