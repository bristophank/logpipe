package renamer

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
	rn := New(nil)
	input := `{"level":"info","msg":"hello"}`
	if got := rn.Apply(input); got != input {
		t.Errorf("expected unchanged line, got %s", got)
	}
}

func TestApply_RenamesField(t *testing.T) {
	rn := New([]Rule{{From: "msg", To: "message"}})
	out := decode(t, rn.Apply(`{"level":"info","msg":"hello"}`))
	if _, ok := out["msg"]; ok {
		t.Error("old key 'msg' should not exist")
	}
	if out["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", out["message"])
	}
}

func TestApply_MissingSourceField(t *testing.T) {
	rn := New([]Rule{{From: "nonexistent", To: "target"}})
	input := `{"level":"warn"}`
	out := decode(t, rn.Apply(input))
	if _, ok := out["target"]; ok {
		t.Error("target key should not be added when source is missing")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	rn := New([]Rule{
		{From: "msg", To: "message"},
		{From: "ts", To: "timestamp"},
	})
	out := decode(t, rn.Apply(`{"ts":"2024-01-01","msg":"hi","level":"debug"}`))
	if out["message"] != "hi" {
		t.Errorf("expected message=hi, got %v", out["message"])
	}
	if out["timestamp"] != "2024-01-01" {
		t.Errorf("expected timestamp=2024-01-01, got %v", out["timestamp"])
	}
	if _, ok := out["msg"]; ok {
		t.Error("old key 'msg' should not exist")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	rn := New([]Rule{{From: "a", To: "b"}})
	input := `not json`
	if got := rn.Apply(input); got != input {
		t.Errorf("expected original line on invalid JSON, got %s", got)
	}
}

func TestApply_EmptyRuleSkipped(t *testing.T) {
	rn := New([]Rule{{From: "", To: "b"}, {From: "a", To: ""}})
	if len(rn.rules) != 0 {
		t.Errorf("expected 0 valid rules, got %d", len(rn.rules))
	}
}
