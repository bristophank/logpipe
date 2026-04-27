package scaler

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return obj
}

func TestApply_NoRules(t *testing.T) {
	s := New(nil)
	out, err := s.Apply(`{"val":50}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `{"val":50}` {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s := New([]Rule{{Field: "val", Min: 0, Max: 100, NewMin: 0, NewMax: 1}})
	out, err := s.Apply(`not-json`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `not-json` {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_BasicRescale(t *testing.T) {
	s := New([]Rule{{Field: "score", Min: 0, Max: 100, NewMin: 0, NewMax: 1}})
	out, err := s.Apply(`{"score":50}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	got, _ := obj["score"].(float64)
	if got != 0.5 {
		t.Errorf("expected 0.5, got %v", got)
	}
}

func TestApply_NegativeRange(t *testing.T) {
	s := New([]Rule{{Field: "temp", Min: -40, Max: 60, NewMin: 0, NewMax: 100}})
	out, err := s.Apply(`{"temp":-40}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	got, _ := obj["temp"].(float64)
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	s := New([]Rule{{Field: "missing", Min: 0, Max: 10, NewMin: 0, NewMax: 1}})
	out, err := s.Apply(`{"other":5}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if _, ok := obj["missing"]; ok {
		t.Errorf("field should not be added")
	}
}

func TestApply_ZeroSpan_ReturnsNewMin(t *testing.T) {
	s := New([]Rule{{Field: "v", Min: 5, Max: 5, NewMin: 3, NewMax: 9}})
	out, err := s.Apply(`{"v":5}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	got, _ := obj["v"].(float64)
	if got != 3 {
		t.Errorf("expected 3 (newMin), got %v", got)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	s := New([]Rule{
		{Field: "a", Min: 0, Max: 10, NewMin: 0, NewMax: 1},
		{Field: "b", Min: 0, Max: 100, NewMin: 0, NewMax: 10},
	})
	out, err := s.Apply(`{"a":5,"b":50}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if a, _ := obj["a"].(float64); a != 0.5 {
		t.Errorf("a: expected 0.5, got %v", a)
	}
	if b, _ := obj["b"].(float64); b != 5 {
		t.Errorf("b: expected 5, got %v", b)
	}
}
