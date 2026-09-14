package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

type PayoutAdjustmentRepo struct {
	q *sqlc.Queries
}

func NewPayoutAdjustmentRepo(db *DB) *PayoutAdjustmentRepo {
	return &PayoutAdjustmentRepo{q: sqlc.New(db.Pool)}
}

func newPayoutAdjustmentRepo(q *sqlc.Queries) *PayoutAdjustmentRepo {
	return &PayoutAdjustmentRepo{q: q}
}

func (r *PayoutAdjustmentRepo) Create(ctx context.Context, a payout.Adjustment) (payout.Adjustment, error) {
	row, err := r.q.CreatePayoutAdjustment(ctx, sqlc.CreatePayoutAdjustmentParams{
		ID:              a.ID,
		PayoutID:        a.PayoutID,
		OrderID:         a.OrderID,
		Type:            sqlc.AdjustmentType(a.Type),
		GrossDelta:      fromDecimal(a.GrossDelta),
		CommissionDelta: fromDecimal(a.CommissionDelta),
		NetDelta:        fromDecimal(a.NetDelta),
		Reason:          a.Reason,
	})
	if err != nil {
		return payout.Adjustment{}, errs.Wrap(err, "adjustment repo: create")
	}
	return toAdjustment(row), nil
}

func (r *PayoutAdjustmentRepo) GetByID(ctx context.Context, id uuid.UUID) (payout.Adjustment, error) {
	row, err := r.q.GetPayoutAdjustmentByID(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return payout.Adjustment{}, payout.ErrNotFound
		}
		return payout.Adjustment{}, errs.Wrap(err, "adjustment repo: getByID")
	}
	return toAdjustment(row), nil
}

func (r *PayoutAdjustmentRepo) ListByPayoutID(ctx context.Context, payoutID uuid.UUID) ([]payout.Adjustment, error) {
	rows, err := r.q.ListPayoutAdjustmentsByPayoutID(ctx, payoutID)
	if err != nil {
		return nil, errs.Wrap(err, "adjustment repo: listByPayoutID")
	}

	result := make([]payout.Adjustment, 0, len(rows))
	for _, row := range rows {
		result = append(result, toAdjustment(row))
	}
	return result, nil
}

func (r *PayoutAdjustmentRepo) List(ctx context.Context, p payout.AdjustmentListParams) ([]payout.Adjustment, int64, error) {
	var typ *sqlc.AdjustmentType
	if p.Type != nil {
		t := sqlc.AdjustmentType(*p.Type)
		typ = &t
	}

	rows, err := r.q.ListPayoutAdjustments(ctx, sqlc.ListPayoutAdjustmentsParams{
		PayoutID:    p.PayoutID,
		OrderID:     p.OrderID,
		Type:        typ,
		CreatedFrom: p.CreatedFrom,
		CreatedTo:   p.CreatedTo,
		SortBy:      p.SortBy,
		SortDir:     p.SortDir,
		Limit:       p.Limit,
		Offset:      p.Offset,
	})
	if err != nil {
		return nil, 0, errs.Wrap(err, "adjustment repo: list")
	}

	total, err := r.q.CountPayoutAdjustments(ctx, sqlc.CountPayoutAdjustmentsParams{
		PayoutID:    p.PayoutID,
		OrderID:     p.OrderID,
		Type:        typ,
		CreatedFrom: p.CreatedFrom,
		CreatedTo:   p.CreatedTo,
	})
	if err != nil {
		return nil, 0, errs.Wrap(err, "adjustment repo: list")
	}

	result := make([]payout.Adjustment, 0, len(rows))
	for _, row := range rows {
		result = append(result, toAdjustment(row))
	}
	return result, total, nil
}

func (r *PayoutAdjustmentRepo) SumByPayoutID(ctx context.Context, payoutID uuid.UUID) (payout.AdjustmentSum, error) {
	row, err := r.q.SumPayoutAdjustments(ctx, payoutID)
	if err != nil {
		return payout.AdjustmentSum{}, errs.Wrap(err, "adjustment repo: sumByPayoutID")
	}
	return payout.AdjustmentSum{
		GrossDelta:      toDecimal(row.GrossDelta),
		CommissionDelta: toDecimal(row.CommissionDelta),
		NetDelta:        toDecimal(row.NetDelta),
	}, nil
}

func toAdjustment(row sqlc.PayoutAdjustment) payout.Adjustment {
	return payout.Adjustment{
		ID:              row.ID,
		PayoutID:        row.PayoutID,
		OrderID:         row.OrderID,
		Type:            payout.AdjustmentType(row.Type),
		GrossDelta:      toDecimal(row.GrossDelta),
		CommissionDelta: toDecimal(row.CommissionDelta),
		NetDelta:        toDecimal(row.NetDelta),
		Reason:          row.Reason,
		CreatedAt:       row.CreatedAt,
	}
}
