package testutil

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewMockDatabase() (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, sqlMock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	return sqlxDB, sqlMock
}

type UnknownDatabaseTransaction struct{}

func (m *UnknownDatabaseTransaction) Rollback() error {
	return nil
}

func (m *UnknownDatabaseTransaction) Commit() error {
	return nil
}
