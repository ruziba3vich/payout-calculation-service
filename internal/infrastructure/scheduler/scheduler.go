package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{cron: cron.New()}
}

// Add registers fn on a standard 5 field cron spec, e.g. "0 0 1 * *".
func (s *Scheduler) Add(spec string, name string, fn func(ctx context.Context) error) error {
	_, err := s.cron.AddFunc(spec, func() {
		log.Printf("scheduler: running %s", name)
		if err := fn(context.Background()); err != nil {
			log.Printf("scheduler: %s failed: %v", name, err)
		}
	})
	return errs.Wrap(err, "scheduler: add "+name)
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop(ctx context.Context) {
	done := s.cron.Stop().Done()
	select {
	case <-done:
	case <-ctx.Done():
	}
}
