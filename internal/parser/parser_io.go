package parser

import (
	"bufio"
	"encoding/json"
	"io"
)

// Stream reads lines from r, parses each one, and writes the resulting JSON
// to w. Unparseable lines are silently skipped.
func Stream(r io.Reader, w io.Writer, format Format) error {
	p := New(format)
	scanner := bufio.NewScanner(r)
	enc := json.NewEncoder(w)
	for scanner.Scan() {
		m, err := p.Parse(scanner.Text())
		if err != nil {
			continue
		}
		if err := enc.Encode(m); err != nil {
			return err
		}
	}
	return scanner.Err()
}
