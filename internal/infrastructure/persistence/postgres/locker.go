package postgres

import (
	"context"
	"hash/fnv"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/infrastructure/persistence/postgres/sqlc"
)

// Locker takes a session level advisory lock. The lock lives on one
// dedicated connection, so it stays held across many transactions until Unlock.
type Locker struct {
	db *DB
}

func NewLocker(db *DB) *Locker {
	return &Locker{db: db}
}

// TryLock returns ok=false when another session already holds the lock.
// The returned release func must be called when ok is true.
func (l *Locker) TryLock(ctx context.Context, name string) (release func(), ok bool, err error) {
	conn, err := l.db.Pool.Acquire(ctx)
	if err != nil {
		return nil, false, errs.Wrap(err, "locker: acquire conn")
	}

	key := lockKey(name)
	q := sqlc.New(conn)

	ok, err = q.TryAdvisoryLock(ctx, key)
	if err != nil {
		conn.Release()
		return nil, false, errs.Wrap(err, "locker: try lock")
	}
	if !ok {
		conn.Release()
		return nil, false, nil
	}

	release = func() {
		_, _ = q.AdvisoryUnlock(context.Background(), key)
		conn.Release()
	}
	return release, true, nil
}

func lockKey(name string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(name))
	return int64(h.Sum64())
}
