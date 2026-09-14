package payout

import (
	"context"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
)

// Repos is the set of repositories bound to one transaction.
type Repos struct {
	Couriers    courier.Repository
	Orders      order.Repository
	Payouts     Repository
	Adjustments AdjustmentRepository
}

// Transactor runs fn inside a single database transaction.
type Transactor interface {
	InTx(ctx context.Context, fn func(r Repos) error) error
}
