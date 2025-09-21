package testutil

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type RepositoryDependency struct {
	MockedDB  *sqlx.DB
	MockedSQL sqlmock.Sqlmock
}

func NewRepositoryDependency() RepositoryDependency {
	sqlxDB, sqlMock := NewMockDatabase()

	return RepositoryDependency{
		MockedDB:  sqlxDB,
		MockedSQL: sqlMock,
	}
}
