package postgres

import (
	"context"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

type Transactor struct {
	db *DB
}

func NewTransactor(db *DB) *Transactor {
	return &Transactor{db: db}
}

func (t *Transactor) InTx(ctx context.Context, fn func(r payout.Repos) error) error {
	return t.db.WithTx(ctx, func(q *sqlc.Queries) error {
		return fn(payout.Repos{
			Couriers:    newCourierRepo(q),
			Orders:      newOrderRepo(q),
			Payouts:     newPayoutRepo(q),
			Adjustments: newPayoutAdjustmentRepo(q),
		})
	})
}
