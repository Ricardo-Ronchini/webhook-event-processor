package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/internal/cache"
)

type CachedDatabase struct {
	db     *Database
	client cache.CacheClient
	key    string
	ttl    time.Duration
}

// CachedRow — resultado de QueryRow com suporte a cache.
type CachedRow struct {
	db     *Database
	client cache.CacheClient
	key    string
	ttl    time.Duration
	sql    string
	args   []any
	err    error
}

func (cd *CachedDatabase) QueryRow(sql string, args ...any) *CachedRow {
	return &CachedRow{
		db:     cd.db,
		client: cd.client,
		key:    cd.key,
		ttl:    cd.ttl,
		sql:    sql,
		args:   args,
	}
}

// Err espelha sql.Row.Err(): retorna erros de pré-execução (ex: client inválido).
// Na maioria dos casos retorna nil; o erro real vem do Scan.
func (r *CachedRow) Err() error { return r.err }

func (r *CachedRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	if r.client == nil {
		return r.db.QueryRow(r.sql, r.args...).Scan(dest...)
	}

	val, err := r.client.Get(context.Background(), r.key)
	if err == nil {
		// cache hit
		var cached []json.RawMessage
		if err := json.Unmarshal([]byte(val), &cached); err == nil && len(cached) == len(dest) {
			for i, raw := range cached {
				if err = json.Unmarshal(raw, dest[i]); err != nil {
					return err
				}
			}

			return nil
		}
	}

	// cache miss ou erro no Get — vai ao banco
	if err := r.db.QueryRow(r.sql, r.args...).Scan(dest...); err != nil {
		return err
	}

	// serializa e armazena
	vals := make([]json.RawMessage, len(dest))

	for i, d := range dest {
		b, merr := json.Marshal(d)
		if merr != nil {
			return nil // não falha por erro de cache
		}
		vals[i] = b
	}

	if b, merr := json.Marshal(vals); merr == nil {
		_ = r.client.Set(context.Background(), r.key, string(b), r.ttl)
	}

	return nil
}

// CachedRows — resultado de Query multi-row com suporte a cache.
type CachedRows struct {
	// modo cache: dados em memória
	buf    [][]json.RawMessage
	pos    int
	bufErr error

	// modo passthrough: sql.Rows real
	sqlRows *sql.Rows
}

func (cd *CachedDatabase) Query(sqlStr string, args ...any) (*CachedRows, error) {
	if cd.client == nil {
		rows, err := cd.db.Query(sqlStr, args...)
		if err != nil {
			return nil, err
		}

		return &CachedRows{sqlRows: rows}, nil
	}

	val, err := cd.client.Get(context.Background(), cd.key)
	if err == nil {
		// cache hit
		var buf [][]json.RawMessage
		if jsonErr := json.Unmarshal([]byte(val), &buf); jsonErr == nil {
			return &CachedRows{buf: buf, pos: -1}, nil
		}
	}

	// cache miss — executa query e lê tudo
	rows, err := cd.db.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var buf [][]json.RawMessage

	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))

		for i := range vals {
			ptrs[i] = &vals[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := make([]json.RawMessage, len(cols))

		for i, v := range vals {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}

			row[i] = b
		}

		buf = append(buf, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if b, merr := json.Marshal(buf); merr == nil {
		_ = cd.client.Set(context.Background(), cd.key, string(b), cd.ttl)
	}

	return &CachedRows{buf: buf, pos: -1}, nil
}

func (r *CachedRows) Next() bool {
	if r.sqlRows != nil {
		return r.sqlRows.Next()
	}

	r.pos++

	return r.pos < len(r.buf)
}

func (r *CachedRows) Scan(dest ...any) error {
	if r.sqlRows != nil {
		return r.sqlRows.Scan(dest...)
	}

	row := r.buf[r.pos]

	for i, raw := range row {
		if i >= len(dest) {
			break
		}

		if err := json.Unmarshal(raw, dest[i]); err != nil {
			return err
		}
	}

	return nil
}

func (r *CachedRows) Close() error {
	if r.sqlRows != nil {
		return r.sqlRows.Close()
	}

	return nil
}

func (r *CachedRows) Err() error {
	if r.sqlRows != nil {
		return r.sqlRows.Err()
	}

	return r.bufErr
}
