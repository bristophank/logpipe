package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Sink struct {
	Name   string `json:"name"`
	Target string `json:"target"` // "stdout", "stderr", or file path
}

type Route struct {
	Sink  string `json:"sink"`
	Rules []Rule `json:"rules"`
}

type Config struct {
	Format     string  `json:"format"`      // json | pretty | text
	SampleRate uint64  `json:"sample_rate"` // 0/1 = no sampling
	Sinks      []Sink  `json:"sinks"`
	Routes     []Route `json:"routes"`
}

// Load reads and validates a JSON config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
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
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.SampleRate == 0 {
		cfg.SampleRate = 1
	}
	return &cfg, nil
}
