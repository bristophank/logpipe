package spliceor

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the spliceor to each non-empty line,
// and writes results to w. Lines that fail to process are passed through.
func Stream(s *Spliceor, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := s.Apply(line)
		if err != nil {
			out = line
		}
		if _, werr := fmt.Fprintln(w, out); werr != nil {
			return werr
		}
	}
	return scanner.Err()
}
