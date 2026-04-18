package throttle

import (
	"testing"
	"time"
)

func TestAllow_NoCooldown(t *testing.T) {
	th := New(0)
	for i := 0; i < 5; i++ {
		if !th.Allow("line") {
			t.Fatal("expected all lines allowed when cooldown is zero")
		}
	}
}

func TestAllow_FirstOccurrence(t *testing.T) {
	th := New(time.Second)
	if !th.Allow("hello") {
		t.Fatal("first occurrence should be allowed")
	}
}

func TestAllow_SuppressedWithinWindow(t *testing.T) {
	th := New(time.Second)
	th.Allow("dup")
	if th.Allow("dup") {
		t.Fatal("duplicate within window should be suppressed")
	}
}

func TestAllow_AllowedAfterWindow(t *testing.T) {
	now := time.Now()
	th := New(time.Second)
	th.now = func() time.Time { return now }
	th.Allow("line")
	th.now = func() time.Time { return now.Add(2 * time.Second) }
	if !th.Allow("line") {
		t.Fatal("line should be allowed after cooldown expires")
	}
}

func TestAllow_DifferentLinesIndependent(t *testing.T) {
	th := New(time.Second)
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("different lines should be independent")
	}
}

func TestReset(t *testing.T) {
	th := New(time.Second)
	th.Allow("x")
	th.Reset()
	if th.Len() != 0 {
		t.Fatal("expected empty after reset")
	}
	if !th.Allow("x") {
		t.Fatal("line should be allowed after reset")
	}
}

func TestLen(t *testing.T) {
	th := New(time.Second)
	th.Allow("a")
	th.Allow("b")
	th.Allow("a")
	if th.Len() != 2 {
		t.Fatalf("expected 2, got %d", th.Len())
	}
}
