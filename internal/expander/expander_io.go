package expander

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads JSON lines from r, applies the expander, and writes results to w.
// Lines that fail to parse are passed through unchanged.
func (e *Expander) Stream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, _ := e.Apply(line)
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
