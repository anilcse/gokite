package model

import "time"

// Rule mirrors the JSON/YAML/db schema and adds runtime‑friendly helpers.
// All numeric fields are float64 to work directly with live prices.
type Rule struct {
	ID        string
	Name      string
	Enabled   bool
	Symbol    string
	Timeframe string // e.g. "1s", "5m"
	Condition Condition
	Entry     EntryParams
	Exit      ExitParams
}

// --- Condition & indicator helpers ---------------------------------------
type Condition struct {
	Type   string  `json:"type"` // crossover, rsi_lt, etc.
	FastMA int     `json:"fast_ma"`
	SlowMA int     `json:"slow_ma"`
	RsiVal float64 `json:"rsi_val"`
}

// Check returns true if the rule's condition is satisfied at this price/time.
// Simplified: only supports SMA crossover for now.
func (c Condition) Check(price float64, ts time.Time) bool {
	switch c.Type {
	case "crossover":
		// TODO: integrate real SMA cache; placeholder always true for demo.
		return true
	case "rsi_lt":
		return false // placeholder
	default:
		return false
	}
}

// --- Entry / Exit params --------------------------------------------------
type EntryParams struct {
	Side      string `json:"side"` // BUY or SELL
	Qty       int    `json:"qty"`
	OrderType string `json:"order_type"`
}

type ExitParams struct {
	StopLossPct   float64 `json:"stop_loss_pct"`
	TakeProfitPct float64 `json:"take_profit_pct"`
}
