package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

func (db *DB) WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) (err error) {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errs.Wrap(err, "postgres: begin tx")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				err = errors.Join(err, rbErr)
			}
		}
	}()

	if err = fn(sqlc.New(tx)); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return errs.Wrap(err, "postgres: commit tx")
	}
	return nil
}
