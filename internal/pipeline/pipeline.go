package pipeline

import (
	"bufio"
	"io"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/metrics"
	"github.com/user/logpipe/internal/router"
)

// Pipeline reads lines from a reader, filters them, and routes matches.
type Pipeline struct {
	filter    *filter.Filter
	router    *router.Router
	metrics   *metrics.Collector
}

// New creates a Pipeline with the given filter, router, and metrics collector.
func New(f *filter.Filter, r *router.Router, m *metrics.Collector) *Pipeline {
	if m == nil {
		m = metrics.New()
	}
	return &Pipeline{filter: f, router: r, metrics: m}
}

// Run reads from src line by line until EOF or error.
func (p *Pipeline) Run(src io.Reader) error {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		p.metrics.IncIn()
		if !p.filter.Match(line) {
			p.metrics.IncDropped()
			continue
		}
		p.metrics.IncMatched()
		if err := p.router.Route(line); err != nil {
			p.metrics.IncRouteError()
		}
	}
	return scanner.Err()
}

// Metrics returns the collector used by this pipeline.
func (p *Pipeline) Metrics() *metrics.Collector {
	return p.metrics
}
