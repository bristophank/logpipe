package sampler

import (
	"testing"
)

func TestAllow_NoSampling(t *testing.T) {
	s := New(1)
	for i := 0; i < 10; i++ {
		if !s.Allow() {
			t.Fatalf("expected all lines to pass with rate=1, failed at i=%d", i)
		}
	}
}

func TestAllow_ZeroRateActsAsOne(t *testing.T) {
	s := New(0)
	if s.Rate() != 1 {
		t.Fatalf("expected rate=1 for zero input, got %d", s.Rate())
	}
	if !s.Allow() {
		t.Fatal("expected first call to pass")
	}
}

func TestAllow_SamplingRate(t *testing.T) {
	s := New(3)
	passed := 0
	for i := 0; i < 9; i++ {
		if s.Allow() {
			passed++
		}
	}
	if passed != 3 {
		t.Fatalf("expected 3 lines to pass out of 9 with rate=3, got %d", passed)
	}
}

func TestAllow_FirstLineAlwaysPasses(t *testing.T) {
	s := New(5)
	if !s.Allow() {
		t.Fatal("expected first line to always pass")
	}
}

func TestAllow_Reset(t *testing.T) {
	s := New(2)
	s.Allow() // 1 — pass
	s.Allow() // 2 — drop
	s.Reset()
	if !s.Allow() {
		t.Fatal("expected first line after reset to pass")
	}
}

func TestRate_ReturnsConfigured(t *testing.T) {
	s := New(7)
	if s.Rate() != 7 {
		t.Fatalf("expected rate 7, got %d", s.Rate())
	}
}
