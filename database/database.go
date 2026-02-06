package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
	once sync.Once
}

func NewDatabase(ctx context.Context, db_host string, db_user string,
	db_pass string, db_name string, db_port string,
	db_ssl string) (*Database, error) {
	d := &Database{}
	err := d.InitDatabase(ctx, db_host, db_user, db_pass, db_name, db_port, db_ssl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return d, nil
}

func (d *Database) InitDatabase(ctx context.Context, db_host string, db_user string,
	db_pass string, db_name string, db_port string, db_ssl string) error {
	var errmsg error
	port := 0
	if val, err := strconv.Atoi(db_port); err == nil {
		port = val
	}
	d.once.Do(func() {

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s pool_max_conns=30 pool_min_conns=5 pool_max_conn_lifetime=1h pool_max_conn_idle_time=30m",
			db_host, db_user,
			db_pass, db_name, port,
			db_ssl,
		)

		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			errmsg = err
			return
		}

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			errmsg = err
			return
		}

		d.pool = pool
	})
	return errmsg
}

func (d *Database) GetConnection(ctx context.Context) *pgxpool.Pool {
	if d.pool == nil {
		log.Fatal("Database not initialized")
	}
	return d.pool
}

func (d *Database) CloseConnection(ctx context.Context) {
	if d.pool == nil {
		log.Fatalf("Database not initialized")
	}
	d.pool.Close()
}

func (d *Database) UsingTransactions(ctx context.Context) (pgx.Tx, error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (d *Database) QueryTransaction(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *Database) ExecuteTransaction(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) error {
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Rollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Commit(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
