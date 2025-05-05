package kite

import (
	"context"
	"fmt"
	"log"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"

	"github.com/anilcse/gokite/internal/config"
)

// Client is a thin convenience wrapper around gokiteconnect that
// * manages the REST client
// * maintains a (re-connectable) WebSocket for live ticks
// * offers helper methods with retry + basic throttling
//
// It is intentionally light so you can mock it in tests.

type Client struct {
	kc  *kiteconnect.Client
	wsc *kiteticker.Ticker
}

// New returns a ready-to-use Client. AccessToken must already be valid.
func New(kcCfg config.Kite) *Client {
	kc := kiteconnect.New(kcCfg.APIKey)
	kc.SetAccessToken(kcCfg.AccessToken)

	return &Client{
		kc: kc,
	}
}

// MarketOrder places an intraday MARKET order (MIS). Retries up to 3x on transient errors.
func (c *Client) MarketOrder(symbol string, qty int, side string) error {
	params := kiteconnect.OrderParams{
		Exchange:        "NFO", // assumes F&O; change as required
		Tradingsymbol:   symbol,
		TransactionType: side,
		OrderType:       "MARKET",
		Quantity:        qty,
		Product:         "MIS",
		Validity:        "DAY",
	}
	for i := 0; i < 3; i++ {
		if _, err := c.kc.PlaceOrder("regular", params); err == nil {
			return nil
		} else {
			log.Printf("order attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("order failed after retries")
}

// StartWS connects to Kite Ticker and streams ticks via cb.
// It automatically reconnects with backoff if the socket drops.
func (c *Client) StartWS(ctx context.Context, apiKey, accessToken string, tokens []uint32, cb func(tick kitemodels.Tick)) error {
	// Create new ticker instance
	c.wsc = kiteticker.New(apiKey, accessToken)

	// Assign callbacks
	c.wsc.OnError(func(err error) {
		log.Printf("ticker error: %v", err)
	})

	c.wsc.OnConnect(func() {
		log.Println("ticker connected")
		if err := c.wsc.Subscribe(tokens); err != nil {
			log.Printf("subscribe error: %v", err)
			return
		}
		if err := c.wsc.SetMode(kiteticker.ModeFull, tokens); err != nil {
			log.Printf("set mode error: %v", err)
			return
		}
	})

	c.wsc.OnClose(func(code int, reason string) {
		log.Printf("ticker closed: %d - %s", code, reason)
	})

	c.wsc.OnReconnect(func(attempt int, delay time.Duration) {
		log.Printf("ticker reconnecting: attempt %d in %v", attempt, delay)
	})

	c.wsc.OnNoReconnect(func(attempt int) {
		log.Printf("ticker max reconnect attempts reached: %d", attempt)
	})

	c.wsc.OnTick(cb)

	// Start the connection
	go c.wsc.Serve()

	return nil
}
