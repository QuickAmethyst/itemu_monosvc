package sql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Tx interface {
	builtin
	Commit() error
	Rollback() error
	Rebind(query string) string
	Updates(ctx context.Context, tableName string, dest interface{}, whereStruct interface{}) (sql.Result, error)
}

type tx struct {
	tx *sqlx.Tx
}

func (t *tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *tx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

func (t *tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return t.tx.QueryxContext(ctx, query, args...)
}

func (t *tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return t.tx.QueryRowxContext(ctx, query, args...)
}

func (t *tx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.GetContext(ctx, dest, query, args...)
}

func (t *tx) Updates(ctx context.Context, tableName string, dest interface{}, whereStruct interface{}) (sql.Result, error) {
	return Updates(ctx, t, tableName, dest, whereStruct)
}

func (t *tx) Commit() error {
	return t.tx.Commit()
}

func (t *tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *tx) Rebind(query string) string {
	return t.tx.Rebind(query)
}

func NewTx(i *sqlx.Tx) Tx {
	return &tx{i}
}
