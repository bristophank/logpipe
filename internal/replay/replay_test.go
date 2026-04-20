package replay

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// readSeeker wraps a strings.Reader so it satisfies io.ReadSeeker.
type readSeeker struct{ *strings.Reader }

func rs(s string) *readSeeker { return &readSeeker{strings.NewReader(s)} }

func TestRun_SinglePass(t *testing.T) {
	src := rs(`{"level":"info"}` + "\n" + `{"level":"warn"}` + "\n")
	var out bytes.Buffer
	r := New(Config{}, src, &out)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestRun_SkipsEmptyLines(t *testing.T) {
	src := rs("line1\n\nline2\n")
	var out bytes.Buffer
	r := New(Config{}, src, &out)
	_ = r.Run()
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 non-empty lines, got %d", len(lines))
	}
}

func TestRun_Loop(t *testing.T) {
	calls := 0
	src := rs("hello\n")

	// Use a counting writer that stops after 3 writes.
	w := &limitWriter{limit: 3, calls: &calls}
	r := New(Config{Loop: true}, src, w)
	_ = r.Run() // returns when writer errors
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRun_RateDelay(t *testing.T) {
	src := rs("a\nb\n")
	var out bytes.Buffer
	r := New(Config{Rate: 100}, src, &out) // 10ms per line
	start := time.Now()
	_ = r.Run()
	elapsed := time.Since(start)
	// 2 lines at 100 lps = ~20ms; allow generous margin for CI
	if elapsed < 10*time.Millisecond {
		t.Fatalf("expected some delay, got %v", elapsed)
	}
}

func TestRun_NoRate_Fast(t *testing.T) {
	src := rs("a\nb\nc\n")
	var out bytes.Buffer
	r := New(Config{Rate: 0}, src, &out)
	start := time.Now()
	_ = r.Run()
	if time.Since(start) > 50*time.Millisecond {
		t.Fatal("expected near-instant replay with no rate limit")
	}
}

// limitWriter returns an error after limit successful writes.
type limitWriter struct {
	limit int
	calls *int
}

func (lw *limitWriter) Write(p []byte) (int, error) {
	*lw.calls++
	if *lw.calls >= lw.limit {
		return 0, bytes.ErrTooLarge
	}
	return len(p), nil
}
