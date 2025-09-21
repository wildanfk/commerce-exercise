package util_test

import (
	"context"
	"database/sql"
	"order-service/internal/testutil"
	"order-service/internal/util"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type DatabaseTransactionTestMock struct {
	sqlMock sqlmock.Sqlmock
	db      *sqlx.DB
}

func NewDatabaseTransactionWithMock(t *testing.T) *DatabaseTransactionTestMock {
	db, sqlMock := testutil.NewMockDatabase()

	return &DatabaseTransactionTestMock{
		sqlMock: sqlMock,
		db:      db,
	}
}

func TestDatabaseTransactionHandler_Begin(t *testing.T) {
	tests := []struct {
		name     string
		assertFn func()
	}{
		{
			name: "Success Begin Transaction With Options",
			assertFn: func() {
				databaseTransactionWithMock := NewDatabaseTransactionWithMock(t)
				databaseTransactionWithMock.sqlMock.ExpectBegin()

				txHandler := util.NewDatabaseTransactionHandler(databaseTransactionWithMock.db)
				tx, err := txHandler.Begin(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})

				assert.NotNil(t, tx)
				assert.Nil(t, err)
			},
		},
		{
			name: "Success Begin Transaction Without Options",
			assertFn: func() {
				databaseTransactionWithMock := NewDatabaseTransactionWithMock(t)
				databaseTransactionWithMock.sqlMock.ExpectBegin()

				txHandler := util.NewDatabaseTransactionHandler(databaseTransactionWithMock.db)
				tx, err := txHandler.Begin(context.Background(), nil)

				assert.NotNil(t, tx)
				assert.Nil(t, err)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn()
		})
	}
}

func TestDatabaseTransaction_GetExecer(t *testing.T) {
	tests := []struct {
		name     string
		assertFn func()
	}{
		{
			name: "Success GetExecer Without Database Transaction",
			assertFn: func() {
				databaseTransactionWithMock := NewDatabaseTransactionWithMock(t)

				db, err := util.GetExecer(databaseTransactionWithMock.db, nil)
				assert.NotNil(t, db)
				assert.Nil(t, err)
			},
		},
		{
			name: "Success GetExecer With Database Transaction",
			assertFn: func() {
				databaseTransactionWithMock := NewDatabaseTransactionWithMock(t)
				databaseTransactionWithMock.sqlMock.ExpectBegin()

				txHandler := util.NewDatabaseTransactionHandler(databaseTransactionWithMock.db)
				tx, _ := txHandler.Begin(context.Background(), nil)

				db, err := util.GetExecer(databaseTransactionWithMock.db, tx)
				assert.NotNil(t, db)
				assert.Nil(t, err)
			},
		},
		{
			name: "Success GetExecer With Unknown Database Transaction",
			assertFn: func() {
				databaseTransactionWithMock := NewDatabaseTransactionWithMock(t)

				db, err := util.GetExecer(databaseTransactionWithMock.db, &testutil.UnknownDatabaseTransaction{})
				assert.Nil(t, db)
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertFn()
		})
	}
}
