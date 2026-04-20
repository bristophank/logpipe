package sequencer

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the sequencer to each non-empty line,
// and writes results to w. Invalid JSON lines are passed through unchanged.
func (s *Sequencer) Stream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, _ := s.Apply(line)
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
