package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_NoLimit(t *testing.T) {
	l := New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow("sink1") {
			t.Fatal("expected all lines allowed when rate=0")
		}
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	l := New(5)
	for i := 0; i < 5; i++ {
		if !l.Allow("sink1") {
			t.Fatalf("expected line %d to be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New(3)
	for i := 0; i < 3; i++ {
		l.Allow("sink1")
	}
	if l.Allow("sink1") {
		t.Fatal("expected 4th line to be denied")
	}
}

func TestAllow_PerSinkIsolation(t *testing.T) {
	l := New(2)
	l.Allow("a")
	l.Allow("a")
	if l.Allow("a") {
		t.Fatal("sink a should be limited")
	}
	if !l.Allow("b") {
		t.Fatal("sink b should still be allowed")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	l := New(1)
	l.Allow("sink1")
	if l.Allow("sink1") {
		t.Fatal("expected denial within window")
	}
	// Force window expiry
	l.mu.Lock()
	l.window["sink1"] = time.Now().Add(-2 * time.Second)
	l.mu.Unlock()
	if !l.Allow("sink1") {
		t.Fatal("expected allow after window reset")
	}
}

func TestReset(t *testing.T) {
	l := New(1)
	l.Allow("sink1")
	l.Reset()
	if !l.Allow("sink1") {
		t.Fatal("expected allow after Reset")
	}
}
