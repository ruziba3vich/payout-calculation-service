package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

type CourierRepo struct {
	q *sqlc.Queries
}

func NewCourierRepo(db *DB) *CourierRepo {
	return &CourierRepo{q: sqlc.New(db.Pool)}
}

func newCourierRepo(q *sqlc.Queries) *CourierRepo {
	return &CourierRepo{q: q}
}

func (r *CourierRepo) Create(ctx context.Context, c courier.Courier) (courier.Courier, error) {
	row, err := r.q.CreateCourier(ctx, sqlc.CreateCourierParams{
		ID:       c.ID,
		FullName: c.FullName,
		Phone:    c.Phone,
		Password: c.PasswordHash,
		HiredAt:  c.HiredAt,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return courier.Courier{}, courier.ErrPhoneTaken
		}
		return courier.Courier{}, errs.Wrap(err, "courier repo: create")
	}
	return toCourier(row), nil
}

func (r *CourierRepo) GetByID(ctx context.Context, id uuid.UUID) (courier.Courier, error) {
	row, err := r.q.GetCourierByID(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return courier.Courier{}, courier.ErrNotFound
		}
		return courier.Courier{}, errs.Wrap(err, "courier repo: getByID")
	}
	return toCourier(row), nil
}

func (r *CourierRepo) GetByPhone(ctx context.Context, phone string) (courier.Courier, error) {
	row, err := r.q.GetCourierByPhone(ctx, phone)
	if err != nil {
		if isNotFound(err) {
			return courier.Courier{}, courier.ErrNotFound
		}
		return courier.Courier{}, errs.Wrap(err, "courier repo: getByPhone")
	}
	return toCourier(row), nil
}

func (r *CourierRepo) Update(ctx context.Context, p courier.UpdateParams) (courier.Courier, error) {
	row, err := r.q.UpdateCourier(ctx, sqlc.UpdateCourierParams{
		ID:       p.ID,
		FullName: p.FullName,
		Phone:    p.Phone,
		Password: p.PasswordHash,
		HiredAt:  p.HiredAt,
		IsActive: p.IsActive,
	})
	if err != nil {
		switch {
		case isNotFound(err):
			return courier.Courier{}, courier.ErrNotFound
		case isUniqueViolation(err):
			return courier.Courier{}, courier.ErrPhoneTaken
		}
		return courier.Courier{}, errs.Wrap(err, "courier repo: update")
	}
	return toCourier(row), nil
}

func (r *CourierRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return errs.Wrap(r.q.SoftDeleteCourier(ctx, id), "courier repo: softDelete")
}

func (r *CourierRepo) List(ctx context.Context, p courier.ListParams) ([]courier.Courier, int64, error) {
	rows, err := r.q.ListCouriers(ctx, sqlc.ListCouriersParams{
		Limit:  p.Limit,
		Offset: p.Offset,
	})
	if err != nil {
		return nil, 0, errs.Wrap(err, "courier repo: list")
	}

	total, err := r.q.CountCouriers(ctx)
	if err != nil {
		return nil, 0, errs.Wrap(err, "courier repo: list")
	}

	result := make([]courier.Courier, 0, len(rows))
	for _, row := range rows {
		result = append(result, toCourier(row))
	}
	return result, total, nil
}

func (r *CourierRepo) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	ok, err := r.q.CourierExists(ctx, id)
	return ok, errs.Wrap(err, "courier repo: exists")
}

func (r *CourierRepo) ListActiveIDs(ctx context.Context) ([]uuid.UUID, error) {
	ids, err := r.q.ListActiveCourierIDs(ctx)
	return ids, errs.Wrap(err, "courier repo: listActiveIDs")
}

func toCourier(row sqlc.Courier) courier.Courier {
	return courier.Courier{
		ID:           row.ID,
		FullName:     row.FullName,
		Phone:        row.Phone,
		PasswordHash: row.Password,
		HiredAt:      row.HiredAt,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
