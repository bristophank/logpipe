// Package replay provides a log line replayer that re-emits lines
// at a controlled rate, optionally looping over the input.
package replay

import (
	"bufio"
	"io"
	"time"
)

// Config controls replay behaviour.
type Config struct {
	// Rate is the number of lines per second to emit. Zero means no throttle.
	Rate int
	// Loop causes the replayer to restart from the beginning when input is exhausted.
	Loop bool
}

// Replayer reads lines from a source and writes them to a sink at a
// controlled rate.
type Replayer struct {
	cfg    Config
	src    io.ReadSeeker
	sink   io.Writer
	delay  time.Duration
}

// New creates a new Replayer. src must support seeking when Loop is true.
func New(cfg Config, src io.ReadSeeker, sink io.Writer) *Replayer {
	var d time.Duration
	if cfg.Rate > 0 {
		d = time.Second / time.Duration(cfg.Rate)
	}
	return &Replayer{cfg: cfg, src: src, sink: sink, delay: d}
}

// Run starts replaying. It blocks until the source is exhausted (or loops
// forever if Loop is true). Callers may cancel via a surrounding context by
// closing the source.
func (r *Replayer) Run() error {
	for {
		if err := r.pass(); err != nil {
			return err
		}
		if !r.cfg.Loop {
			break
		}
		if _, err := r.src.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}
	return nil
}

func (r *Replayer) pass() error {
	scanner := bufio.NewScanner(r.src)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if _, err := r.sink.Write(append(line, '\n')); err != nil {
			return err
		}
		if r.delay > 0 {
			time.Sleep(r.delay)
		}
	}
	return scanner.Err()
}
