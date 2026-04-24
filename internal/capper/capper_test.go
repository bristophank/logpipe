package capper

import (
	"fmt"
	"testing"
	"time"
)

func line(field, value string) string {
	return fmt.Sprintf(`{%q:%q}`, field, value)
}

func TestAllow_NoRules(t *testing.T) {
	c := New(nil)
	if !c.Allow(line("level", "error")) {
		t.Fatal("expected allow with no rules")
	}
}

func TestAllow_WithinCap(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 3, Window: time.Minute}})
	for i := 0; i < 3; i++ {
		if !c.Allow(line("level", "error")) {
			t.Fatalf("expected allow on iteration %d", i)
		}
	}
}

func TestAllow_ExceedsCap(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 2, Window: time.Minute}})
	c.Allow(line("level", "warn"))
	c.Allow(line("level", "warn"))
	if c.Allow(line("level", "warn")) {
		t.Fatal("expected drop after exceeding cap")
	}
}

func TestAllow_DifferentValuesIndependent(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 1, Window: time.Minute}})
	if !c.Allow(line("level", "error")) {
		t.Fatal("error should be allowed")
	}
	if !c.Allow(line("level", "warn")) {
		t.Fatal("warn should be allowed independently")
	}
	if c.Allow(line("level", "error")) {
		t.Fatal("second error should be dropped")
	}
}

func TestAllow_MissingField_Passes(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 1, Window: time.Minute}})
	for i := 0; i < 5; i++ {
		if !c.Allow(`{"msg":"hello"}`) {
			t.Fatal("line without tracked field should always pass")
		}
	}
}

func TestAllow_InvalidJSON_Passes(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 1, Window: time.Minute}})
	if !c.Allow("not-json") {
		t.Fatal("invalid JSON should always pass")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 1, Window: 10 * time.Millisecond}})
	c.Allow(line("level", "error"))
	if c.Allow(line("level", "error")) {
		t.Fatal("expected drop before window reset")
	}
	time.Sleep(20 * time.Millisecond)
	if !c.Allow(line("level", "error")) {
		t.Fatal("expected allow after window reset")
	}
}

func TestAllow_Reset(t *testing.T) {
	c := New([]Rule{{Field: "level", Max: 1, Window: time.Minute}})
	c.Allow(line("level", "error"))
	c.Reset()
	if !c.Allow(line("level", "error")) {
		t.Fatal("expected allow after Reset")
	}
}
