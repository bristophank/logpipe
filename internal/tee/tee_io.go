package tee

import (
	"bufio"
	"io"
)

// Stream reads lines from r and fans each line out to all registered writers
// until r is exhausted or returns an error.
func (t *Tee) Stream(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := append(scanner.Bytes(), '\n')
		if _, err := t.Write(line); err != nil {
			return err
		}
	}
	return scanner.Err()
}
