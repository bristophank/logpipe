package dispatcher

import (
	"bufio"
	"io"
)

// Stream reads lines from r, dispatches each through d, and returns the total
// number of lines dispatched (non-empty). Any write errors are returned
// immediately.
func Stream(r io.Reader, d *Dispatcher) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		sink, err := d.Dispatch(line)
		if err != nil {
			return count, err
		}
		if sink != "" {
			count++
		}
	}
	return count, scanner.Err()
}
