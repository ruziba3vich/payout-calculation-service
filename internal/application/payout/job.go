package payout

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
)

// MonthlyJob calculates the previous month's payout for every active courier.
type MonthlyJob struct {
	payouts  *Service
	couriers courier.Repository
	locker   payout.Locker
}

func NewMonthlyJob(payouts *Service, couriers courier.Repository, locker payout.Locker) *MonthlyJob {
	return &MonthlyJob{payouts: payouts, couriers: couriers, locker: locker}
}

type JobResult struct {
	Period   string `json:"period"`
	Skipped  bool   `json:"skipped"`
	Reason   string `json:"reason,omitempty"`
	Total    int    `json:"total"`
	Created  int    `json:"created"`
	Existing int    `json:"existing"`
	Failed   int    `json:"failed"`
}

// RunPrevious runs the job for the month before now.
func (j *MonthlyJob) RunPrevious(ctx context.Context) (JobResult, error) {
	prev := payout.PeriodOf(time.Now()).AddDate(0, -1, 0)
	return j.Run(ctx, prev)
}

// Run calculates payouts for one period. Two instances calling this at the
// same time: the advisory lock lets only one through, the other returns Skipped.
// Each courier is its own transaction, so one failure does not stop the rest.
func (j *MonthlyJob) Run(ctx context.Context, period time.Time) (JobResult, error) {
	period = payout.PeriodOf(period)
	res := JobResult{Period: period.Format("2006-01")}

	release, ok, err := j.locker.TryLock(ctx, "payout-job:"+res.Period)
	if err != nil {
		return res, err
	}
	if !ok {
		res.Skipped = true
		res.Reason = "another instance is running this period"
		log.Printf("payout job %s: skipped, lock held", res.Period)
		return res, nil
	}
	defer release()

	ids, err := j.couriers.ListActiveIDs(ctx)
	if err != nil {
		return res, err
	}
	res.Total = len(ids)

	for _, id := range ids {
		_, err := j.payouts.Calculate(ctx, id, period)
		switch {
		case err == nil:
			res.Created++
		case errors.Is(err, payout.ErrAlreadyExists):
			res.Existing++
		default:
			res.Failed++
			log.Printf("payout job %s: courier %s: %v", res.Period, id, err)
		}
	}

	log.Printf("payout job %s: total=%d created=%d existing=%d failed=%d",
		res.Period, res.Total, res.Created, res.Existing, res.Failed)
	return res, nil
}
