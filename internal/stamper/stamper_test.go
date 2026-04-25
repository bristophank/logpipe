package stamper

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	s := New(nil)
	input := `{"msg":"hello"}`
	if got := s.Apply(input); got != input {
		t.Fatalf("expected passthrough, got %s", got)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s := New([]Rule{{Field: "n", Start: 0, Step: 1}})
	input := `not-json`
	if got := s.Apply(input); got != input {
		t.Fatalf("expected passthrough on invalid JSON, got %s", got)
	}
}

func TestApply_DefaultStartAndStep(t *testing.T) {
	s := New([]Rule{{Field: "seq"}})
	m := decode(t, s.Apply(`{"msg":"a"}`))
	if int(m["seq"].(float64)) != 0 {
		t.Fatalf("expected seq=0, got %v", m["seq"])
	}
	m = decode(t, s.Apply(`{"msg":"b"}`))
	if int(m["seq"].(float64)) != 1 {
		t.Fatalf("expected seq=1, got %v", m["seq"])
	}
}

func TestApply_CustomStartAndStep(t *testing.T) {
	s := New([]Rule{{Field: "n", Start: 10, Step: 5}})
	m := decode(t, s.Apply(`{"x":1}`))
	if int(m["n"].(float64)) != 10 {
		t.Fatalf("expected n=10, got %v", m["n"])
	}
	m = decode(t, s.Apply(`{"x":2}`))
	if int(m["n"].(float64)) != 15 {
		t.Fatalf("expected n=15, got %v", m["n"])
	}
}

func TestApply_StringFormat(t *testing.T) {
	s := New([]Rule{{Field: "id", Start: 1, Step: 1, Format: "string"}})
	m := decode(t, s.Apply(`{"msg":"hi"}`))
	if m["id"] != "1" {
		t.Fatalf("expected id=\"1\", got %v", m["id"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	s := New([]Rule{
		{Field: "a", Start: 0, Step: 1},
		{Field: "b", Start: 100, Step: 10},
	})
	m := decode(t, s.Apply(`{"msg":"x"}`))
	if int(m["a"].(float64)) != 0 || int(m["b"].(float64)) != 100 {
		t.Fatalf("unexpected values: a=%v b=%v", m["a"], m["b"])
	}
	m = decode(t, s.Apply(`{"msg":"y"}`))
	if int(m["a"].(float64)) != 1 || int(m["b"].(float64)) != 110 {
		t.Fatalf("unexpected values after step: a=%v b=%v", m["a"], m["b"])
	}
}

func TestApply_Reset(t *testing.T) {
	s := New([]Rule{{Field: "n", Start: 5, Step: 1}})
	s.Apply(`{"x":1}`)
	s.Apply(`{"x":2}`)
	s.Reset()
	m := decode(t, s.Apply(`{"x":3}`))
	if int(m["n"].(float64)) != 5 {
		t.Fatalf("expected reset to start=5, got %v", m["n"])
	}
}
