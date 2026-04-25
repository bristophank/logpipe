package clamper

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	c := New(nil)
	in := `{"value":200}`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Fatalf("expected unchanged, got %s", out)
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	c := New([]Rule{{Field: "score", Min: 0, Max: 100}})
	in := `{"score":50}`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Fatalf("expected unchanged, got %s", out)
	}
}

func TestApply_BelowMin_ClampsToMin(t *testing.T) {
	c := New([]Rule{{Field: "score", Min: 0, Max: 100}})
	out, err := c.Apply(`{"score":-10}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["score"].(float64) != 0 {
		t.Fatalf("expected 0, got %v", m["score"])
	}
}

func TestApply_AboveMax_ClampsToMax(t *testing.T) {
	c := New([]Rule{{Field: "score", Min: 0, Max: 100}})
	out, err := c.Apply(`{"score":150}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["score"].(float64) != 100 {
		t.Fatalf("expected 100, got %v", m["score"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c := New([]Rule{{Field: "score", Min: 0, Max: 100}})
	in := `{"level":"info"}`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Fatalf("expected unchanged, got %s", out)
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	c := New([]Rule{{Field: "score", Min: 0, Max: 100}})
	in := `not-json`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Fatalf("expected passthrough, got %s", out)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	c := New([]Rule{
		{Field: "score", Min: 0, Max: 100},
		{Field: "latency", Min: 1, Max: 5000},
	})
	out, err := c.Apply(`{"score":999,"latency":-3}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["score"].(float64) != 100 {
		t.Fatalf("score: expected 100, got %v", m["score"])
	}
	if m["latency"].(float64) != 1 {
		t.Fatalf("latency: expected 1, got %v", m["latency"])
	}
}
