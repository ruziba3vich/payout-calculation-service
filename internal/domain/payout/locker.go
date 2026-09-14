package payout

import "context"

// Locker gives a cluster wide lock so only one service instance runs a job.
type Locker interface {
	TryLock(ctx context.Context, name string) (release func(), ok bool, err error)
}
