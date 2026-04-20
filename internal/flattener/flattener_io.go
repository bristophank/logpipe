package flattener

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the flattener to each line,
// and writes the result to w. Empty lines are skipped.
func Stream(f *Flattener, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out := f.Apply(line)
		if _, err := fmt.Fprintln(w, out); err != nil {
			return err
		}
	}
	return scanner.Err()
}
