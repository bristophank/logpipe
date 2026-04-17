package dedup

import (
	"testing"
	"time"
)

func TestAllow_NoWindow(t *testing.T) {
	d := New(0)
	if !d.Allow(`{"msg":"hello"}`) {
		t.Fatal("expected allow when window is zero")
	}
	if !d.Allow(`{"msg":"hello"}`) {
		t.Fatal("expected allow on repeat when window is zero")
	}
}

func TestAllow_UniqueLines(t *testing.T) {
	d := New(5 * time.Second)
	if !d.Allow(`{"msg":"a"}`) {
		t.Fatal("expected first unique line to pass")
	}
	if !d.Allow(`{"msg":"b"}`) {
		t.Fatal("expected second unique line to pass")
	}
}

func TestAllow_DuplicateSuppressed(t *testing.T) {
	d := New(5 * time.Second)
	line := `{"msg":"dup"}`
	if !d.Allow(line) {
		t.Fatal("expected first occurrence to pass")
	}
	if d.Allow(line) {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	base := time.Now()
	d := New(2 * time.Second)
	d.now = func() time.Time { return base }

	line := `{"msg":"expire"}`
	if !d.Allow(line) {
		t.Fatal("expected first pass")
	}
	if d.Allow(line) {
		t.Fatal("expected suppression within window")
	}

	// Advance time past window
	d.now = func() time.Time { return base.Add(3 * time.Second) }
	if !d.Allow(line) {
		t.Fatal("expected pass after window expiry")
	}
}

func TestAllow_Reset(t *testing.T) {
	d := New(10 * time.Second)
	line := `{"msg":"reset"}`
	d.Allow(line)
	d.Reset()
	if !d.Allow(line) {
		t.Fatal("expected pass after reset")
	}
}

func TestAllow_IsolatedKeys(t *testing.T) {
	d := New(10 * time.Second)
	for i := 0; i < 5; i++ {
		line := `{"i":` + string(rune('0'+i)) + `}`
		if !d.Allow(line) {
			t.Fatalf("expected unique line %d to pass", i)
		}
	}
}
