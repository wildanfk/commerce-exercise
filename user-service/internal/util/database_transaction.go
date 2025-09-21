package util

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=mock/database_transaction.go -package=mock -source=database_transaction.go

type DatabaseTransaction interface {
	Rollback() error
	Commit() error
}

type DatabaseTransactionHandler interface {
	Begin(ctx context.Context, opts *sql.TxOptions) (DatabaseTransaction, error)
}

type databaseTransactionHandler struct {
	db *sqlx.DB
}

func NewDatabaseTransactionHandler(db *sqlx.DB) *databaseTransactionHandler {
	return &databaseTransactionHandler{db: db}
}

func (dth *databaseTransactionHandler) Begin(ctx context.Context, opts *sql.TxOptions) (DatabaseTransaction, error) {
	return dth.db.BeginTxx(ctx, opts)
}

type DatabaseTransactionAction interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error)
}

func GetExecer(db *sqlx.DB, tx DatabaseTransaction) (DatabaseTransactionAction, error) {
	if tx == nil {
		return db, nil
	}

	sqltx, ok := tx.(*sqlx.Tx)
	if !ok {
		return nil, errors.New("database transaction not supported")
	}

	return sqltx, nil
}
