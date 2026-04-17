package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/logpipe/internal/transformer"
)

// SinkConfig describes an output destination.
type SinkConfig struct {
	Name string `json:"name"`
	Type string `json:"type"` // stdout, file
	Path string `json:"path"` // for file sinks
}

// RouteConfig maps a sink name to filter rules.
type RouteConfig struct {
	Sink  string            `json:"sink"`
	Rules []map[string]string `json:"rules"`
}

// Config is the top-level configuration structure.
type Config struct {
	Sinks        []SinkConfig         `json:"sinks"`
	Routes       []RouteConfig        `json:"routes"`
	Transformers []transformer.Rule   `json:"transformers"`
}

// Load reads and validates a JSON config file.
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
		if !sinkNames[r.Sink] {
			return nil, fmt.Errorf("config: route references unknown sink %q", r.Sink)
		}
	}
	for _, tr := range cfg.Transformers {
		switch tr.Op {
		case "set", "delete", "rename", "uppercase", "lowercase":
		default:
			return nil, fmt.Errorf("config: unknown transformer op %q", tr.Op)
		}
	}
	return &cfg, nil
}
