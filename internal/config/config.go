package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/logpipe/internal/redactor"
)

// SinkConfig describes an output destination.
type SinkConfig struct {
	Name string `json:"name"`
	Path string `json:"path"` // "-" for stdout
}

// RouteConfig maps a filter set to a sink.
type RouteConfig struct {
	Sink    string            `json:"sink"`
	Filters map[string]string `json:"filters"`
}

// Config is the top-level configuration structure.
type Config struct {
	Sinks     []SinkConfig      `json:"sinks"`
	Routes    []RouteConfig     `json:"routes"`
	Redact    []redactor.Rule   `json:"redact"`
	Format    string            `json:"format"`
	RateLimit int               `json:"rate_limit"`
	SampleN   int               `json:"sample_n"`
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
		if s.Name == "" {
			return nil, fmt.Errorf("config: sink missing name")
		}
		sinkNames[s.Name] = true
	}
	for _, r := range cfg.Routes {
		if !sinkNames[r.Sink] {
			return nil, fmt.Errorf("config: route references unknown sink %q", r.Sink)
		}
	}
	return &cfg, nil
}
