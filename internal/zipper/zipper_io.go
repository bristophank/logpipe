package zipper

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the zipper to each, and writes results to w.
// Lines that fail to parse are passed through unchanged.
func (z *Zipper) Stream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := z.Apply(line)
		if err != nil {
			out = line
		}
		if _, werr := io.WriteString(w, out+"\n"); werr != nil {
			return werr
		}
	}
	return scanner.Err()
}
