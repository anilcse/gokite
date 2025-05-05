package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anilcse/gokite/internal/config"
	"github.com/anilcse/gokite/internal/engine"
	"github.com/anilcse/gokite/internal/kite"
	"github.com/anilcse/gokite/internal/scheduler"
	"github.com/anilcse/gokite/internal/store"
)

func main() {
	// 1. Load application configuration (app.yaml)
	cfg, err := config.Load("configs/app.yaml")
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	// 2. Connect to postgres (or whichever DSN you provided)
	db, err := store.Connect(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	// 3. Spin up Kite client (REST + WS)
	kc := kite.New(cfg.Kite)

	// 4. Create rule-evaluation engine
	eng := engine.New(db, kc)

	// 5. Scheduler: cron jobs + rule hot-reload
	sch := scheduler.New(eng, db)

	// 6. Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go sch.Start(ctx)
	// Convert config.Kite to engine.Config
	wsConfig := engine.Config{
		APIKey:      cfg.Kite.APIKey,
		AccessToken: cfg.Kite.AccessToken,
		Instruments: cfg.Kite.Instruments,
	}
	go eng.StartWS(ctx, wsConfig) // subscribes to live ticks

	// 7. Wait for SIGINT/SIGTERM to stop cleanly
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("shutting down…")
	cancel()

	// allow goroutines to flush
	time.Sleep(2 * time.Second)
	log.Println("bye")
}
