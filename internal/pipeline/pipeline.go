package pipeline

import (
	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/formatter"
	"github.com/yourorg/logpipe/internal/metrics"
	"github.com/yourorg/logpipe/internal/ratelimit"
	"github.com/yourorg/logpipe/internal/router"
)

// Pipeline wires filter → formatter → rate-limiter → router.
type Pipeline struct {
	filter    *filter.Filter
	router    *router.Router
	formatter *formatter.Formatter
	metrics   *metrics.Collector
	limiter   *ratelimit.Limiter
}

// New constructs a Pipeline from its dependencies.
func New(
	f *filter.Filter,
	r *router.Router,
	fmt *formatter.Formatter,
	m *metrics.Collector,
	l *ratelimit.Limiter,
) *Pipeline {
	return &Pipeline{filter: f, router: r, formatter: fmt, metrics: m, limiter: l}
}

// Process handles a single input line.
func (p *Pipeline) Process(line string) {
	if line == "" {
		return
	}
	p.metrics.IncProcessed()

	if !p.filter.Match(line) {
		p.metrics.IncDropped()
		return
	}

	formatted := p.formatter.Format(line)

	sinks := p.router.Sinks()
	for _, name := range sinks {
		if !p.limiter.Allow(name) {
			p.metrics.IncDropped()
			continue
		}
		p.router.Write(name, formatted)
		p.metrics.IncRouted()
	}
}
