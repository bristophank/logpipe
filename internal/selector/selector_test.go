package selector

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("failed to decode output: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	s := New(nil)
	input := `{"level":"info","msg":"hello","host":"web-1"}`
	if got := s.Apply(input); got != input {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestApply_KeepFields(t *testing.T) {
	s := New([]Rule{{Fields: []string{"level", "msg"}, Mode: "keep"}})
	input := `{"level":"info","msg":"hello","host":"web-1"}`
	out := decode(t, s.Apply(input))
	if _, ok := out["level"]; !ok {
		t.Error("expected 'level' to be kept")
	}
	if _, ok := out["msg"]; !ok {
		t.Error("expected 'msg' to be kept")
	}
	if _, ok := out["host"]; ok {
		t.Error("expected 'host' to be dropped")
	}
}

func TestApply_DropFields(t *testing.T) {
	s := New([]Rule{{Fields: []string{"host", "pid"}, Mode: "drop"}})
	input := `{"level":"error","msg":"oops","host":"db-1","pid":42}`
	out := decode(t, s.Apply(input))
	if _, ok := out["host"]; ok {
		t.Error("expected 'host' to be dropped")
	}
	if _, ok := out["pid"]; ok {
		t.Error("expected 'pid' to be dropped")
	}
	if _, ok := out["level"]; !ok {
		t.Error("expected 'level' to remain")
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	s := New([]Rule{{Fields: []string{"level"}, Mode: "keep"}})
	input := `not-json`
	if got := s.Apply(input); got != input {
		t.Errorf("expected passthrough for invalid JSON, got %s", got)
	}
}

func TestApply_EmptyLine(t *testing.T) {
	s := New([]Rule{{Fields: []string{"level"}, Mode: "keep"}})
	if got := s.Apply(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestApply_KeepMissingField(t *testing.T) {
	s := New([]Rule{{Fields: []string{"level", "nonexistent"}, Mode: "keep"}})
	input := `{"level":"warn","msg":"test"}`
	out := decode(t, s.Apply(input))
	if _, ok := out["level"]; !ok {
		t.Error("expected 'level' to be kept")
	}
	if _, ok := out["msg"]; ok {
		t.Error("expected 'msg' to be dropped")
	}
}
