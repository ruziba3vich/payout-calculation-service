package payout

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/cache"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
)

// Cache keys:
//   payout:{id}                         -> payout with adjustments, TTL PayoutTTL
//   payouts:courier:{courierID}:ver     -> version counter, no TTL
//   payouts:courier:{courierID}:{ver}:{filters} -> list page, TTL PayoutListTTL
//
// Invalidation: a payout changes (new adjustment, status change) -> DEL payout:{id}
// and INCR the courier version so every cached list page for that courier is orphaned.
// Orphaned pages expire on their own TTL.

type Cache struct {
	store   cache.Store
	itemTTL time.Duration
	listTTL time.Duration
}

func NewCache(store cache.Store, itemTTL, listTTL time.Duration) *Cache {
	return &Cache{store: store, itemTTL: itemTTL, listTTL: listTTL}
}

type cachedList struct {
	Items []payout.Payout
	Total int64
}

func payoutKey(id uuid.UUID) string {
	return "payout:" + id.String()
}

func courierVerKey(courierID uuid.UUID) string {
	return "payouts:courier:" + courierID.String() + ":ver"
}

func (c *Cache) getPayout(ctx context.Context, id uuid.UUID) (PayoutWithAdjustments, bool) {
	if c == nil {
		return PayoutWithAdjustments{}, false
	}
	var v PayoutWithAdjustments
	ok, err := c.store.Get(ctx, payoutKey(id), &v)
	if err != nil {
		log.Println("cache:", err)
		return PayoutWithAdjustments{}, false
	}
	return v, ok
}

func (c *Cache) setPayout(ctx context.Context, v PayoutWithAdjustments) {
	if c == nil {
		return
	}
	if err := c.store.Set(ctx, payoutKey(v.Payout.ID), v, c.itemTTL); err != nil {
		log.Println("cache:", err)
	}
}

func (c *Cache) listKey(ctx context.Context, courierID uuid.UUID, p payout.ListParams) (string, bool) {
	ver, err := c.store.GetInt(ctx, courierVerKey(courierID))
	if err != nil {
		log.Println("cache:", err)
		return "", false
	}

	status := ""
	if p.Status != nil {
		status = string(*p.Status)
	}
	from, to := "", ""
	if p.PeriodFrom != nil {
		from = p.PeriodFrom.Format("2006-01")
	}
	if p.PeriodTo != nil {
		to = p.PeriodTo.Format("2006-01")
	}

	return fmt.Sprintf("payouts:courier:%s:%d:%s:%s:%s:%s:%s:%d:%d",
		courierID, ver, status, from, to, p.SortBy, p.SortDir, p.Limit, p.Offset), true
}

func (c *Cache) getList(ctx context.Context, courierID uuid.UUID, p payout.ListParams) ([]payout.Payout, int64, bool) {
	if c == nil {
		return nil, 0, false
	}
	key, ok := c.listKey(ctx, courierID, p)
	if !ok {
		return nil, 0, false
	}
	var v cachedList
	found, err := c.store.Get(ctx, key, &v)
	if err != nil {
		log.Println("cache:", err)
		return nil, 0, false
	}
	if !found {
		return nil, 0, false
	}
	return v.Items, v.Total, true
}

func (c *Cache) setList(ctx context.Context, courierID uuid.UUID, p payout.ListParams, items []payout.Payout, total int64) {
	if c == nil {
		return
	}
	key, ok := c.listKey(ctx, courierID, p)
	if !ok {
		return
	}
	if err := c.store.Set(ctx, key, cachedList{Items: items, Total: total}, c.listTTL); err != nil {
		log.Println("cache:", err)
	}
}

// invalidate drops the payout item and bumps the courier list version.
func (c *Cache) invalidate(ctx context.Context, payoutID, courierID uuid.UUID) {
	if c == nil {
		return
	}
	if err := c.store.Del(ctx, payoutKey(payoutID)); err != nil {
		log.Println("cache:", err)
	}
	if _, err := c.store.Incr(ctx, courierVerKey(courierID)); err != nil {
		log.Println("cache:", err)
	}
}
