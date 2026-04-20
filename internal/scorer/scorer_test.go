package scorer

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	s := New(Config{})
	out := s.Apply(`{"msg":"hello"}`)
	m := decode(t, out)
	if m["score"] != float64(0) {
		t.Fatalf("expected score 0, got %v", m["score"])
	}
}

func TestApply_MatchingRule(t *testing.T) {
	s := New(Config{
		Rules: []Rule{{Field: "level", Contains: "error", Score: 10}},
	})
	out := s.Apply(`{"level":"error","msg":"boom"}`)
	m := decode(t, out)
	if m["score"] != float64(10) {
		t.Fatalf("expected score 10, got %v", m["score"])
	}
}

func TestApply_NonMatchingRule(t *testing.T) {
	s := New(Config{
		Rules: []Rule{{Field: "level", Contains: "error", Score: 10}},
	})
	out := s.Apply(`{"level":"info","msg":"ok"}`)
	m := decode(t, out)
	if m["score"] != float64(0) {
		t.Fatalf("expected score 0, got %v", m["score"])
	}
}

func TestApply_MultipleRules_Cumulative(t *testing.T) {
	s := New(Config{
		Rules: []Rule{
			{Field: "level", Contains: "error", Score: 10},
			{Field: "msg", Contains: "timeout", Score: 5},
		},
	})
	out := s.Apply(`{"level":"error","msg":"timeout occurred"}`)
	m := decode(t, out)
	if m["score"] != float64(15) {
		t.Fatalf("expected score 15, got %v", m["score"])
	}
}

func TestApply_CustomOutputKey(t *testing.T) {
	s := New(Config{
		OutputKey: "priority",
		Rules:     []Rule{{Field: "level", Contains: "warn", Score: 3}},
	})
	out := s.Apply(`{"level":"warn"}`)
	m := decode(t, out)
	if _, ok := m["priority"]; !ok {
		t.Fatal("expected key 'priority' in output")
	}
	if m["priority"] != float64(3) {
		t.Fatalf("expected priority 3, got %v", m["priority"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	s := New(Config{})
	input := "not json"
	if got := s.Apply(input); got != input {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestApply_EmptyLine_Passthrough(t *testing.T) {
	s := New(Config{})
	if got := s.Apply(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestApply_MissingField_ZeroContribution(t *testing.T) {
	s := New(Config{
		Rules: []Rule{{Field: "nonexistent", Contains: "x", Score: 99}},
	})
	out := s.Apply(`{"level":"info"}`)
	m := decode(t, out)
	if m["score"] != float64(0) {
		t.Fatalf("expected score 0, got %v", m["score"])
	}
}
