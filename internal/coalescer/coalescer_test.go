package coalescer

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	c := New(nil)
	input := `{"msg":"hello"}`
	if got := c.Apply(input); got != input {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestApply_FirstCandidateWins(t *testing.T) {
	c := New([]Rule{{Target: "message", Candidates: []string{"msg", "text", "body"}}})
	m := decode(t, c.Apply(`{"msg":"hello","body":"world"}`))
	if m["message"] != "hello" {
		t.Errorf("expected 'hello', got %v", m["message"])
	}
}

func TestApply_SkipsEmptyString(t *testing.T) {
	c := New([]Rule{{Target: "message", Candidates: []string{"msg", "text"}}})
	m := decode(t, c.Apply(`{"msg":"","text":"fallback"}`))
	if m["message"] != "fallback" {
		t.Errorf("expected 'fallback', got %v", m["message"])
	}
}

func TestApply_DefaultWhenNoCandidateFound(t *testing.T) {
	c := New([]Rule{{Target: "level", Candidates: []string{"lvl", "severity"}, Default: "info"}})
	m := decode(t, c.Apply(`{"msg":"hi"}`))
	if m["level"] != "info" {
		t.Errorf("expected default 'info', got %v", m["level"])
	}
}

func TestApply_NoDefaultNoMatch_TargetAbsent(t *testing.T) {
	c := New([]Rule{{Target: "level", Candidates: []string{"lvl"}}})
	m := decode(t, c.Apply(`{"msg":"hi"}`))
	if _, ok := m["level"]; ok {
		t.Error("expected target to be absent")
	}
}

func TestApply_DeleteSrc_RemovesCandidates(t *testing.T) {
	c := New([]Rule{{Target: "message", Candidates: []string{"msg", "text"}, DeleteSrc: true}})
	m := decode(t, c.Apply(`{"msg":"hello","text":"world"}`))
	if m["message"] != "hello" {
		t.Errorf("expected 'hello', got %v", m["message"])
	}
	if _, ok := m["msg"]; ok {
		t.Error("expected 'msg' to be deleted")
	}
	if _, ok := m["text"]; ok {
		t.Error("expected 'text' to be deleted")
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	c := New([]Rule{{Target: "x", Candidates: []string{"a"}}})
	input := `not-json`
	if got := c.Apply(input); got != input {
		t.Errorf("expected passthrough for invalid JSON, got %s", got)
	}
}

func TestApply_EmptyLine_Passthrough(t *testing.T) {
	c := New([]Rule{{Target: "x", Candidates: []string{"a"}}})
	if got := c.Apply(""); got != "" {
		t.Errorf("expected empty passthrough, got %q", got)
	}
}
