package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type DB interface {
	BeginTx(ctx context.Context) (pgx.Tx, DB, error)
	Close() error
	Connection() Connection
}

type db struct {
	pool *pgxpool.Pool
	conn Connection
}

func NewDB(postgresURL string) (DB, error) {
	connConfig, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database")

	return &db{
		pool: pool,
		conn: pool,
	}, nil
}

func (d *db) BeginTx(ctx context.Context) (pgx.Tx, DB, error) {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}

	return tx, &db{
		pool: d.pool,
		conn: tx.Conn(),
	}, nil
}

func (d *db) Close() error {
	log.Println("Closing database connection")
	d.pool.Close()
	log.Println("Database connection closed")

	return nil
}

func (d *db) Connection() Connection {
	return d.conn
}
