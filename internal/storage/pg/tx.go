package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

var ErrNoTx = errors.New("в контексте нет транзакции")

// Transactor выполняет функцию внутри транзакции pgx, пряча саму транзакцию
// в контексте — так слой приложения не зависит от драйвера БД.
type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{pool: pool}
}

func (t *Transactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("закоммитить транзакцию: %w", err)
	}

	return nil
}

// TxFromContext достаёт транзакцию, открытую Transactor.WithinTx.
func TxFromContext(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return nil, ErrNoTx
	}

	return tx, nil
}
