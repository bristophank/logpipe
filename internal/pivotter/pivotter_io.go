package pivotter

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads JSON lines from r, applies the Pivotter to each line,
// and writes the result to w. Empty lines are skipped.
func Stream(p *Pivotter, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out := p.Apply(line)
		if _, err := fmt.Fprintln(w, out); err != nil {
			return err
		}
	}
	return scanner.Err()
}
