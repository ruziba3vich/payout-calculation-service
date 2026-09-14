package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

type PayoutRepo struct {
	q *sqlc.Queries
}

func NewPayoutRepo(db *DB) *PayoutRepo {
	return &PayoutRepo{q: sqlc.New(db.Pool)}
}

func newPayoutRepo(q *sqlc.Queries) *PayoutRepo {
	return &PayoutRepo{q: q}
}

func (r *PayoutRepo) Create(ctx context.Context, p payout.Payout) (payout.Payout, error) {
	row, err := r.q.CreatePayout(ctx, sqlc.CreatePayoutParams{
		ID:               p.ID,
		CourierID:        p.CourierID,
		Period:           p.Period,
		DeliveredCount:   p.DeliveredCount,
		GrossAmount:      fromDecimal(p.GrossAmount),
		CommissionRate:   fromDecimal(p.CommissionRate),
		CommissionAmount: fromDecimal(p.CommissionAmount),
		NetAmount:        fromDecimal(p.NetAmount),
	})
	if err != nil {
		if isUniqueViolation(err) {
			return payout.Payout{}, payout.ErrAlreadyExists
		}
		return payout.Payout{}, errs.Wrap(err, "payout repo: create")
	}
	return toPayout(row), nil
}

func (r *PayoutRepo) GetByID(ctx context.Context, id uuid.UUID) (payout.Payout, error) {
	row, err := r.q.GetPayoutByID(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return payout.Payout{}, payout.ErrNotFound
		}
		return payout.Payout{}, errs.Wrap(err, "payout repo: getByID")
	}
	return toPayout(row), nil
}

func (r *PayoutRepo) GetByCourierPeriod(ctx context.Context, courierID uuid.UUID, period time.Time) (payout.Payout, error) {
	row, err := r.q.GetPayoutByCourierPeriod(ctx, sqlc.GetPayoutByCourierPeriodParams{
		CourierID: courierID,
		Period:    period,
	})
	if err != nil {
		if isNotFound(err) {
			return payout.Payout{}, payout.ErrNotFound
		}
		return payout.Payout{}, errs.Wrap(err, "payout repo: getByCourierPeriod")
	}
	return toPayout(row), nil
}

func (r *PayoutRepo) GetByCourierPeriodForUpdate(ctx context.Context, courierID uuid.UUID, period time.Time) (payout.Payout, error) {
	row, err := r.q.GetPayoutByCourierPeriodForUpdate(ctx, sqlc.GetPayoutByCourierPeriodForUpdateParams{
		CourierID: courierID,
		Period:    period,
	})
	if err != nil {
		if isNotFound(err) {
			return payout.Payout{}, payout.ErrNotFound
		}
		return payout.Payout{}, errs.Wrap(err, "payout repo: getByCourierPeriodForUpdate")
	}
	return toPayout(row), nil
}

func (r *PayoutRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status payout.Status) (payout.Payout, error) {
	row, err := r.q.UpdatePayoutStatus(ctx, sqlc.UpdatePayoutStatusParams{
		ID:     id,
		Status: sqlc.PayoutStatus(status),
	})
	if err != nil {
		if isNotFound(err) {
			return payout.Payout{}, payout.ErrNotFound
		}
		return payout.Payout{}, errs.Wrap(err, "payout repo: updateStatus")
	}
	return toPayout(row), nil
}

func (r *PayoutRepo) List(ctx context.Context, p payout.ListParams) ([]payout.Payout, int64, error) {
	var status *sqlc.PayoutStatus
	if p.Status != nil {
		s := sqlc.PayoutStatus(*p.Status)
		status = &s
	}

	rows, err := r.q.ListPayouts(ctx, sqlc.ListPayoutsParams{
		CourierID:  p.CourierID,
		Status:     status,
		PeriodFrom: p.PeriodFrom,
		PeriodTo:   p.PeriodTo,
		SortBy:     p.SortBy,
		SortDir:    p.SortDir,
		Limit:      p.Limit,
		Offset:     p.Offset,
	})
	if err != nil {
		return nil, 0, errs.Wrap(err, "payout repo: list")
	}

	total, err := r.q.CountPayouts(ctx, sqlc.CountPayoutsParams{
		CourierID:  p.CourierID,
		Status:     status,
		PeriodFrom: p.PeriodFrom,
		PeriodTo:   p.PeriodTo,
	})
	if err != nil {
		return nil, 0, errs.Wrap(err, "payout repo: list")
	}

	result := make([]payout.Payout, 0, len(rows))
	for _, row := range rows {
		result = append(result, toPayout(row))
	}
	return result, total, nil
}

func toPayout(row sqlc.Payout) payout.Payout {
	return payout.Payout{
		ID:               row.ID,
		CourierID:        row.CourierID,
		Period:           row.Period,
		DeliveredCount:   row.DeliveredCount,
		GrossAmount:      toDecimal(row.GrossAmount),
		CommissionRate:   toDecimal(row.CommissionRate),
		CommissionAmount: toDecimal(row.CommissionAmount),
		NetAmount:        toDecimal(row.NetAmount),
		Status:           payout.Status(row.Status),
		CalculatedAt:     row.CalculatedAt,
		PaidAt:           row.PaidAt,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}
