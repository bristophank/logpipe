package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// SinkConfig describes an output destination.
type SinkConfig struct {
	Name   string `json:"name"`
	Type   string `json:"type"` // "stdout", "file"
	Target string `json:"target,omitempty"`
}

// RouteConfig maps a filter rule set to a list of sink names.
type RouteConfig struct {
	Sinks  []string          `json:"sinks"`
	Filter map[string]string `json:"filter,omitempty"`
}

// Config is the top-level configuration structure.
type Config struct {
	Sinks   []SinkConfig `json:"sinks"`
	Routes  []RouteConfig `json:"routes"`
	RateLimit int         `json:"rate_limit,omitempty"` // lines/sec per sink, 0 = unlimited
	Format  string       `json:"format,omitempty"`    // json|pretty|text
}

// Load reads and validates a JSON config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}
	sinkNames := make(map[string]bool, len(cfg.Sinks))
	for _, s := range cfg.Sinks {
		sinkNames[s.Name] = true
	}
	for _, r := range cfg.Routes {
		for _, sn := range r.Sinks {
			if !sinkNames[sn] {
				return nil, fmt.Errorf("config: route references unknown sink %q", sn)
			}
		}
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	return &cfg, nil
}
