package splitter

import "io"

// WriterMap is a convenience alias used when constructing a Splitter from
// named io.Writer values.
type WriterMap map[string]io.Writer

// NewFromWriters constructs a Splitter directly from an io.WriterMap.
func NewFromWriters(rules []Rule, sinks WriterMap, fallback io.Writer) *Splitter {
	w := make(map[string]io.Writer, len(sinks))
	for k, v := range sinks {
		w[k] = v
	}
	return New(rules, w, fallback)
}
