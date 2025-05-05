package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// App holds the full application configuration loaded from configs/app.yaml
// Environment variables (e.g., DB_DSN) can override individual fields when present.
//
// Example YAML (configs/app.yaml):
// ---
// env: development
// kite:
//   api_key: "KITE_API_KEY"
//   api_secret: "KITE_API_SECRET"
//   access_token: "KITE_ACCESS_TOKEN" # refreshed daily via external script
//   instruments: ["NFO:NIFTY24MAYFUT", "NFO:BANKNIFTY24MAYFUT"]
// db:
//   dsn: "postgres://user:pass@localhost:5432/algo_kite?sslmode=disable"
//   max_idle: 5
//   max_open: 10
//   log_queries: false
//
// You can mount different YAMLs per environment (app.dev.yaml, app.prod.yaml, …)
// and pass the path via APP_CONFIG env var or CLI flag.

type App struct {
	Env  string `yaml:"env"`
	Kite Kite   `yaml:"kite"`
	DB   DB     `yaml:"db"`
}

type Kite struct {
	APIKey      string   `yaml:"api_key"`
	APISecret   string   `yaml:"api_secret"`
	AccessToken string   `yaml:"access_token"`
	Instruments []string `yaml:"instruments"`
}

type DB struct {
	DSN        string `yaml:"dsn"`
	MaxIdle    int    `yaml:"max_idle"`
	MaxOpen    int    `yaml:"max_open"`
	LogQueries bool   `yaml:"log_queries"`
}

// Load reads the YAML file at path and unmarshals it into App.
// On error it returns non‑nil.
func Load(path string) (App, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return App{}, fmt.Errorf("read config: %w", err)
	}
	var cfg App
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return App{}, fmt.Errorf("yaml parse: %w", err)
	}
	// TODO(optional): overlay env var overrides here
	return cfg, nil
}
