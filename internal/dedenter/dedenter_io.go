package dedenter

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads newline-delimited JSON lines from r, applies the Dedenter to
// each non-empty line, and writes the result to w. It returns the first write
// error encountered, or nil on EOF.
func Stream(r io.Reader, w io.Writer, d *Dedenter) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out := d.Apply(line)
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
