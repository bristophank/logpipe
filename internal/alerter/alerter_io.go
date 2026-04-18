package alerter

import (
	"bufio"
	"io"
)

// Stream reads lines from r, checks each against the Alerter rules,
// and writes passing lines to out. Alerts are written to alertOut.
func (a *Alerter) Stream(r io.Reader, out io.Writer, alertOut io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if _, err := io.WriteString(out, line+"\n"); err != nil {
			return err
		}
		if err := a.Check(line, alertOut); err != nil {
			return err
		}
	}
	return scanner.Err()
}
