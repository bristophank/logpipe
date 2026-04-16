package pipeline

import (
	"bufio"
	"io"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/router"
)

// Config holds pipeline configuration.
type Config struct {
	SinkName string
	Rules    []filter.Rule
}

// Pipeline ties a filter and router together, reading from a source.
type Pipeline struct {
	filter *filter.Filter
	router *router.Router
	cfg    Config
}

// New creates a Pipeline with the given filter rules and router.
func New(cfg Config, r *router.Router) (*Pipeline, error) {
	f, err := filter.New(cfg.Rules)
	if err != nil {
		return nil, err
	}
	return &Pipeline{filter: f, router: r, cfg: cfg}, nil
}

// Run reads lines from src, applies filter rules, and routes matching lines.
func (p *Pipeline) Run(src io.Reader) error {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		ok, err := p.filter.Match(line)
		if err != nil || !ok {
			continue
		}
		if p.cfg.SinkName != "" {
			if err := p.router.Route(line, p.cfg.SinkName); err != nil {
				return err
			}
		} else {
			if err := p.router.Route(line); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
