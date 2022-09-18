package sql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Result = sql.Result

type DB interface {
	Stats() sql.DBStats
	Close() error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PingContext(ctx context.Context) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type db struct {
	db *sqlx.DB
}

func (d *db) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args)
}

func (d *db) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

func (d *db) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.GetContext(ctx, dest, query, args...)
}

func (d *db) Rebind(query string) string {
	return d.db.Rebind(query)
}

func (d *db) PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return d.db.PreparexContext(ctx, query)
}

func (d *db) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return d.db.QueryRowxContext(ctx, query, args...)
}

func (d *db) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.db.QueryxContext(ctx, query, args)
}

func (d *db) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *db) Stats() sql.DBStats {
	return d.db.Stats()
}

func (d *db) Close() error {
	return d.db.Close()
}

func (d *db) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.SelectContext(ctx, dest, query, args...)
}

func NewDB(i *sqlx.DB) DB {
	return &db{i}
}
