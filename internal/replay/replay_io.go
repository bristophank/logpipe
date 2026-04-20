package replay

import (
	"io"
	"os"
)

// Stream opens the file at path and replays its contents to sink according
// to cfg. It is a convenience wrapper around New and Run suitable for use
// from cmd/logpipe.
func Stream(cfg Config, path string, sink io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return New(cfg, f, sink).Run()
}
