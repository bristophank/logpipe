package selector

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads newline-delimited JSON from r, applies field selection to each
// line, and writes the result to w. Empty lines are skipped. Lines that cannot
// be processed are passed through unchanged so the stream is never silently
// dropped.
func Stream(r io.Reader, w io.Writer, s *Selector) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		out, err := s.Apply(line)
		if err != nil {
			// Pass invalid JSON through unmodified.
			out = line
		}

		if _, werr := io.WriteString(w, out+"\n"); werr != nil {
			return werr
		}
	}
	return scanner.Err()
}
