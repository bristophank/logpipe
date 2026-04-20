package joiner

import (
	"bufio"
	"fmt"
	"io"
)

// Stream reads lines from r, applies the joiner to each, and writes results to w.
// Empty lines are skipped. Lines that fail to marshal are passed through unchanged.
func Stream(j *Joiner, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		out, _ := j.Apply(line)
		if _, err := fmt.Fprintln(w, out); err != nil {
			return err
		}
	}
	return scanner.Err()
}
