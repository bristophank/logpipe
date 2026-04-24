package unwrapper

import (
	"bufio"
	"io"
	"strings"
)

// Stream reads lines from r, applies the unwrapper to each, and writes results to w.
func Stream(u *Unwrapper, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := u.Apply(line)
		if err != nil {
			continue
		}
		if _, err := io.WriteString(w, out+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
