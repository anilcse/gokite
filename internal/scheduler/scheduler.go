package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/anilcse/gokite/internal/engine"
	"github.com/anilcse/gokite/internal/store"
)

// Scheduler periodically reloads strategy rules and runs housekeeping jobs.
// It is kept minimal; you can add cron expressions or a full job registry later.

type Scheduler struct {
	eng *engine.Engine
	db  *store.DB
}

func New(eng *engine.Engine, db *store.DB) *Scheduler {
	return &Scheduler{eng: eng, db: db}
}

func (s *Scheduler) Start(ctx context.Context) {
	// initial load
	if err := s.eng.LoadRules(); err != nil {
		log.Printf("rule load: %v", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := s.eng.LoadRules(); err != nil {
				log.Printf("rule reload: %v", err)
			}
		}
	}
}
