package classifier

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the classifier to each non-empty line,
// and writes the result to w. Malformed JSON lines are passed through unchanged.
func (c *Classifier) Stream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := c.Apply(line)
		if err != nil {
			out = line
		}
		if _, werr := io.WriteString(w, out+"\n"); werr != nil {
			return werr
		}
	}
	return scanner.Err()
}
