package counter

import (
	"bufio"
	"fmt"
	"io"
)

// Stream reads newline-delimited JSON lines from r, passes each through Add,
// and writes every non-empty line to w unchanged.
// Add errors (invalid JSON) are silently ignored so the stream is never interrupted.
func (c *Counter) Stream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		_ = c.Add(line)
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return err
		}
	}
	return scanner.Err()
}
