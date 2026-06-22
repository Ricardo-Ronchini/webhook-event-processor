package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/cache"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase() *Database {
	return &Database{conn: Connect()}
}

func (d *Database) WithCache(client cache.CacheClient, key string, ttl time.Duration) *CachedDatabase {
	return &CachedDatabase{db: d, client: client, key: key, ttl: ttl}
}

func (d *Database) Close() error {
	return d.conn.Close()
}

func (d *Database) ConfigurePool(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) {
	d.conn.SetMaxOpenConns(maxOpen)
	d.conn.SetMaxIdleConns(maxIdle)
	d.conn.SetConnMaxLifetime(maxLifetime)
	d.conn.SetConnMaxIdleTime(maxIdleTime)
}

func (d *Database) Begin() (*sql.Tx, error) {
	return d.conn.Begin()
}

func (d *Database) Exec(query string, args ...any) (sql.Result, error) {
	return d.conn.Exec(query, args...)
}

func (d *Database) ExecWithContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.conn.ExecContext(ctx, query, args...)
}

func (d *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return d.conn.Query(query, args...)
}

func (d *Database) QueryWithContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.conn.QueryContext(ctx, query, args...)
}

func (d *Database) QueryRow(query string, args ...any) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

func (d *Database) QueryRowWithContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.conn.QueryRowContext(ctx, query, args...)
}
