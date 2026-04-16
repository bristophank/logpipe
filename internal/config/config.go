package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Rule defines a single filter rule.
type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value,omitempty"`
}

// Sink defines an output destination.
type Sink struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Route maps a sink name to a set of filter rules.
type Route struct {
	Sink  string `json:"sink"`
	Rules []Rule `json:"rules"`
}

// Config is the top-level configuration structure.
type Config struct {
	Sinks  []Sink  `json:"sinks"`
	Routes []Route `json:"routes"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	sinkNames := make(map[string]struct{}, len(c.Sinks))
	for _, s := range c.Sinks {
		if s.Name == "" {
			return fmt.Errorf("sink missing name")
		}
		if s.Path == "" {
			return fmt.Errorf("sink %q missing path", s.Name)
		}
		sinkNames[s.Name] = struct{}{}
	}
	for _, r := range c.Routes {
		if _, ok := sinkNames[r.Sink]; !ok {
			return fmt.Errorf("route references unknown sink %q", r.Sink)
		}
	}
	return nil
}
