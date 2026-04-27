package shifter

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the shifter to each, and writes results
// to w. Empty lines are skipped. Non-JSON lines are passed through unchanged.
func Stream(s *Shifter, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := s.Apply(line)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
