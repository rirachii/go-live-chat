package db

import (
	"context"
	"database/sql"
	"fmt"

	// _ "github.com/lib/pq"
	pgx "github.com/jackc/pgx/v5"
)

type Database struct {
	db *pgx.Conn
}

// uses pgx driver
func ConnectDatabase() (*Database, error) {

	const (
		driver   = "pgx"
		username = "root"
		password = "password"
		host     = "localhost:5432"
		options  = "sslmode=disable"
	)
	var DB_URL = fmt.Sprintf(
		"postgres://%s:%s@%s/go-chat?%s",
		username, password, host, options,
	)

	db, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close(ctx context.Context) {
	d.db.Close(ctx)
}

func (d *Database) DB() *pgx.Conn {
	return d.db
}

type SQL_DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type PGX_DBTX interface {
}
