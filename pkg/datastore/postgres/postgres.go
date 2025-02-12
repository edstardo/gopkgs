package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PostgresDB struct {
	Pool Querier
}

type TxFunc func(ctx context.Context, q Querier) error

func New(url string) (*PostgresDB, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		Pool: pool,
	}, nil
}

func (p *PostgresDB) Close() {
	p.Pool.(*pgxpool.Pool).Close()
}

func (p *PostgresDB) Tx(ctx context.Context, fn TxFunc) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}