package store

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/anilcse/gokite/internal/model"
)

// DB wraps sqlx.DB to expose typed helper methods.

type DB struct {
	*sqlx.DB
}

// Connect opens a PostgreSQL connection pool.
func Connect(dsn string) (*DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// ActiveRules returns all enabled strategy_rules rows mapped to model.Rule.
func (db *DB) ActiveRules() ([]model.Rule, error) {
	var rows []struct {
		ID        string `db:"id"`
		Name      string `db:"name"`
		Enabled   bool   `db:"enabled"`
		Symbol    string `db:"symbol"`
		Timeframe string `db:"timeframe"`
		Defn      string `db:"defn"`
	}
	if err := db.Select(&rows, `SELECT id,name,enabled,symbol,timeframe,defn FROM strategy_rules WHERE enabled = true`); err != nil {
		return nil, err
	}
	out := make([]model.Rule, 0, len(rows))
	for _, r := range rows {
		rule := model.Rule{
			ID:        r.ID,
			Name:      r.Name,
			Enabled:   r.Enabled,
			Symbol:    r.Symbol,
			Timeframe: r.Timeframe,
		}
		// naive JSON decode into Condition/Entry/Exit omitted for brevity
		out = append(out, rule)
	}
	return out, nil
}

// InstrumentTokens maps trading symbols to instrument tokens via a lookup
// table you should keep synced from Kite instruments dump.
func (db *DB) InstrumentTokens(symbols []string) ([]uint32, error) {
	type row struct {
		Token  uint32 `db:"instrument_token"`
		Symbol string `db:"tradingsymbol"`
	}
	query, args, _ := sqlx.In(`SELECT instrument_token, tradingsymbol FROM instruments WHERE tradingsymbol IN (?)`, symbols)
	var rs []row
	if err := db.Select(&rs, db.Rebind(query), args...); err != nil {
		return nil, err
	}
	tokens := make([]uint32, 0, len(rs))
	for _, r := range rs {
		tokens = append(tokens, r.Token)
	}
	return tokens, nil
}
