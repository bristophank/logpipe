package router

import (
	"io"
	"sync"
)

// Sink represents a named output destination.
type Sink struct {
	Name   string
	Writer io.Writer
}

// Router routes log lines to one or more sinks.
type Router struct {
	mu    sync.RWMutex
	sinks map[string]*Sink
}

// New creates a new Router.
func New() *Router {
	return &Router{
		sinks: make(map[string]*Sink),
	}
}

// AddSink registers a named sink.
func (r *Router) AddSink(name string, w io.Writer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sinks[name] = &Sink{Name: name, Writer: w}
}

// RemoveSink unregisters a named sink.
func (r *Router) RemoveSink(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sinks, name)
}

// Route writes line to the specified sinks. If no names are given, writes to all sinks.
func (r *Router) Route(line []byte, names ...string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(names) == 0 {
		for _, s := range r.sinks {
			if err := write(s.Writer, line); err != nil {
				return err
			}
		}
	}

	for _, name := range names {
		s, ok := r.sinks[name]
		if !ok {
			continue
		}
		if err := write(s.Writer, line); err != nil {
			return err
		}
	}
	return nil
}

func write(w io.Writer, line []byte) error {
	if len(line) == 0 {
		return nil
	}
	_, err := w.Write(append(line, '\n'))
	return err
}
