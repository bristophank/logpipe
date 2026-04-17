package transformer

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	tr := New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := tr.Apply(line); got != line {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestApply_SetField(t *testing.T) {
	tr := New([]Rule{{Field: "env", Op: "set", Value: "production"}})
	m := decode(t, tr.Apply(`{"msg":"hi"}`))
	if m["env"] != "production" {
		t.Errorf("expected env=production, got %v", m["env"])
	}
}

func TestApply_DeleteField(t *testing.T) {
	tr := New([]Rule{{Field: "secret", Op: "delete"}})
	m := decode(t, tr.Apply(`{"msg":"hi","secret":"abc"}`))
	if _, ok := m["secret"]; ok {
		t.Error("expected secret to be deleted")
	}
}

func TestApply_RenameField(t *testing.T) {
	tr := New([]Rule{{Field: "msg", Op: "rename", Value: "message"}})
	m := decode(t, tr.Apply(`{"msg":"hello"}`))
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
	if _, ok := m["msg"]; ok {
		t.Error("old field should be gone")
	}
}

func TestApply_Uppercase(t *testing.T) {
	tr := New([]Rule{{Field: "level", Op: "uppercase"}})
	m := decode(t, tr.Apply(`{"level":"info"}`))
	if m["level"] != "INFO" {
		t.Errorf("expected INFO, got %v", m["level"])
	}
}

func TestApply_Lowercase(t *testing.T) {
	tr := New([]Rule{{Field: "level", Op: "lowercase"}})
	m := decode(t, tr.Apply(`{"level":"ERROR"}`))
	if m["level"] != "error" {
		t.Errorf("expected error, got %v", m["level"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	tr := New([]Rule{{Field: "x", Op: "set", Value: "y"}})
	line := "not json"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected passthrough on bad json")
	}
}
