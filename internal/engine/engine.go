package engine

import (
	"context"
	"log"
	"sync"
	"time"

	kitemodels "github.com/zerodha/gokiteconnect/v4/models"

	"github.com/anilcse/gokite/internal/kite"
	"github.com/anilcse/gokite/internal/model"
	"github.com/anilcse/gokite/internal/store"
)

// Engine coordinates rule evaluation against live market data and delegates
// order execution to the Kite client.
//
// * Pulls enabled rules from DB (hot-reloadable)
// * Maintains an in-memory map keyed by (symbol, timeframe)
// * Evaluates each tick/candle and fires orders via kite.Client
//
// The Engine is stateless across restarts because it can reconstruct all
// indicators from persisted candles/positional data.

type Engine struct {
	db   *store.DB
	kite *kite.Client

	mu      sync.RWMutex
	rules   map[string][]model.Rule // key = symbol+timeframe
	updates chan struct{}           // notify of rule reload
}

func New(db *store.DB, kc *kite.Client) *Engine {
	e := &Engine{
		db:      db,
		kite:    kc,
		rules:   make(map[string][]model.Rule),
		updates: make(chan struct{}, 1),
	}
	return e
}

// Config holds Kite Connect API configuration.
type Config struct {
	APIKey      string
	AccessToken string
	Instruments []string
}

// StartWS is a convenience proxy so main.go can call eng.StartWS(ctx,…).
// It retrieves the list of instrument tokens from config and forwards ticks
// into OnTick for evaluation.
func (e *Engine) StartWS(ctx context.Context, kcCfg Config) {
	tokens, err := e.db.InstrumentTokens(kcCfg.Instruments)
	if err != nil {
		log.Fatalf("token lookup: %v", err)
	}
	go e.kite.StartWS(ctx, kcCfg.APIKey, kcCfg.AccessToken, tokens, e.OnTick)
}

// LoadRules fetches all enabled rules from DB into memory.
func (e *Engine) LoadRules() error {
	rules, err := e.db.ActiveRules()
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = make(map[string][]model.Rule)
	for _, r := range rules {
		key := r.Symbol + r.Timeframe
		e.rules[key] = append(e.rules[key], r)
	}
	select {
	case e.updates <- struct{}{}:
	default:
	}
	return nil
}

// OnTick is passed as callback to Kite WS. It batches ticks by 1-second window
// and forwards aggregated OHLCV candles into evaluate().
func (e *Engine) OnTick(tick kitemodels.Tick) {
	// For brevity, we treat each tick as a candle close (1-second timeframe).
	// Get the symbol from the instrument token mapping
	symbol := "UNKNOWN" // TODO: Implement instrument token to symbol mapping
	e.evaluate(symbol, tick.LastPrice, time.Now())
}

func (e *Engine) evaluate(symbol string, price float64, ts time.Time) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	key := symbol + "1s" // hard-coded timeframe in this snippet
	for _, r := range e.rules[key] {
		if r.Enabled && r.Condition.Check(price, ts) {
			side := "BUY"
			if r.Entry.Side == "SELL" {
				side = "SELL"
			}
			if err := e.kite.MarketOrder(symbol, r.Entry.Qty, side); err != nil {
				log.Printf("order error: %v", err)
			}
		}
	}
}
