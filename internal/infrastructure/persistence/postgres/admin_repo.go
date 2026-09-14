package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/admin"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

type AdminRepo struct {
	q *sqlc.Queries
}

func NewAdminRepo(db *DB) *AdminRepo {
	return &AdminRepo{q: sqlc.New(db.Pool)}
}

func (r *AdminRepo) Create(ctx context.Context, a admin.Admin) (admin.Admin, error) {
	row, err := r.q.CreateAdmin(ctx, sqlc.CreateAdminParams{
		ID:       a.ID,
		FullName: a.FullName,
		Username: a.Username,
		Password: a.PasswordHash,
	})
	if err != nil {
		return admin.Admin{}, err
	}
	return toAdmin(row), nil
}

func (r *AdminRepo) GetByID(ctx context.Context, id uuid.UUID) (admin.Admin, error) {
	row, err := r.q.GetAdminByID(ctx, id)
	if err != nil {
		return admin.Admin{}, err
	}
	return toAdmin(row), nil
}

func (r *AdminRepo) GetByUsername(ctx context.Context, username string) (admin.Admin, error) {
	row, err := r.q.GetAdminByUsername(ctx, username)
	if err != nil {
		return admin.Admin{}, err
	}
	return toAdmin(row), nil
}

func (r *AdminRepo) Update(ctx context.Context, p admin.UpdateParams) (admin.Admin, error) {
	row, err := r.q.UpdateAdmin(ctx, sqlc.UpdateAdminParams{
		ID:       p.ID,
		FullName: p.FullName,
		Username: p.Username,
		Password: p.PasswordHash,
	})
	if err != nil {
		return admin.Admin{}, err
	}
	return toAdmin(row), nil
}

func (r *AdminRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteAdmin(ctx, id)
}

func (r *AdminRepo) List(ctx context.Context, p admin.ListParams) ([]admin.Admin, int64, error) {
	rows, err := r.q.ListAdmins(ctx, sqlc.ListAdminsParams{
		Limit:  p.Limit,
		Offset: p.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountAdmins(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]admin.Admin, 0, len(rows))
	for _, row := range rows {
		result = append(result, toAdmin(row))
	}
	return result, total, nil
}

func toAdmin(row sqlc.Administration) admin.Admin {
	return admin.Admin{
		ID:           row.ID,
		FullName:     row.FullName,
		Username:     row.Username,
		PasswordHash: row.Password,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
