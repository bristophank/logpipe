package tagger

import (
	"bufio"
	"io"
)

// Stream reads lines from r, applies tagging rules, and writes results to w.
// Lines that fail to parse are passed through unchanged.
func (t *Tagger) Stream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		out, _ := t.Apply(line)
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
