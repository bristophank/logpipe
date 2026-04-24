package sanitizer

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies sanitization rules to each line,
// and writes the result to w. Empty lines are skipped.
func Stream(r io.Reader, w io.Writer, s *Sanitizer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out := s.Apply(line)
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
