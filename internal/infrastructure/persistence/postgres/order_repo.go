package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

type OrderRepo struct {
	q *sqlc.Queries
}

func NewOrderRepo(db *DB) *OrderRepo {
	return &OrderRepo{q: sqlc.New(db.Pool)}
}

func newOrderRepo(q *sqlc.Queries) *OrderRepo {
	return &OrderRepo{q: q}
}

func (r *OrderRepo) Create(ctx context.Context, o order.Order) (order.Order, error) {
	row, err := r.q.CreateOrder(ctx, sqlc.CreateOrderParams{
		ID:        o.ID,
		CourierID: o.CourierID,
		Amount:    fromDecimal(o.Amount),
	})
	if err != nil {
		return order.Order{}, err
	}
	return toOrder(row), nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id uuid.UUID) (order.Order, error) {
	row, err := r.q.GetOrderByID(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return order.Order{}, order.ErrNotFound
		}
		return order.Order{}, err
	}
	return toOrder(row), nil
}

func (r *OrderRepo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (order.Order, error) {
	row, err := r.q.GetOrderByIDForUpdate(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return order.Order{}, order.ErrNotFound
		}
		return order.Order{}, err
	}
	return toOrder(row), nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status order.Status) (order.Order, error) {
	row, err := r.q.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     id,
		Status: sqlc.OrderStatus(status),
	})
	if err != nil {
		if isNotFound(err) {
			return order.Order{}, order.ErrNotFound
		}
		return order.Order{}, err
	}
	return toOrder(row), nil
}

func (r *OrderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteOrder(ctx, id)
}

func (r *OrderRepo) List(ctx context.Context, p order.ListParams) ([]order.Order, int64, error) {
	var status *sqlc.OrderStatus
	if p.Status != nil {
		s := sqlc.OrderStatus(*p.Status)
		status = &s
	}

	rows, err := r.q.ListOrders(ctx, sqlc.ListOrdersParams{
		CourierID:     p.CourierID,
		Status:        status,
		DeliveredFrom: p.DeliveredFrom,
		DeliveredTo:   p.DeliveredTo,
		CreatedFrom:   p.CreatedFrom,
		CreatedTo:     p.CreatedTo,
		AmountMin:     fromDecimalPtr(p.AmountMin),
		AmountMax:     fromDecimalPtr(p.AmountMax),
		SortBy:        p.SortBy,
		SortDir:       p.SortDir,
		Limit:         p.Limit,
		Offset:        p.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountOrders(ctx, sqlc.CountOrdersParams{
		CourierID:     p.CourierID,
		Status:        status,
		DeliveredFrom: p.DeliveredFrom,
		DeliveredTo:   p.DeliveredTo,
		CreatedFrom:   p.CreatedFrom,
		CreatedTo:     p.CreatedTo,
		AmountMin:     fromDecimalPtr(p.AmountMin),
		AmountMax:     fromDecimalPtr(p.AmountMax),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]order.Order, 0, len(rows))
	for _, row := range rows {
		result = append(result, toOrder(row))
	}
	return result, total, nil
}

func (r *OrderRepo) DeliveredStats(ctx context.Context, courierID uuid.UUID, start, end time.Time) (int64, decimal.Decimal, error) {
	row, err := r.q.GetDeliveredStatsForPeriod(ctx, sqlc.GetDeliveredStatsForPeriodParams{
		CourierID:   courierID,
		PeriodStart: start,
		PeriodEnd:   end,
	})
	if err != nil {
		return 0, decimal.Zero, err
	}
	return row.DeliveredCount, toDecimal(row.DeliveredTotal), nil
}

func toOrder(row sqlc.Order) order.Order {
	return order.Order{
		ID:          row.ID,
		CourierID:   row.CourierID,
		Amount:      toDecimal(row.Amount),
		Status:      order.Status(row.Status),
		DeliveredAt: row.DeliveredAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
