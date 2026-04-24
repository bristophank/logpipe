package sanitizer

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
	s := New(nil)
	input := `{"msg":"<b>hello</b>"}`
	if got := s.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_StripHTML(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "strip_html"}})
	out := s.Apply(`{"msg":"<b>hello</b> world"}`)
	m := decode(t, out)
	if m["msg"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", m["msg"])
	}
}

func TestApply_Trim(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "trim"}})
	out := s.Apply(`{"msg":"  hello  "}`)
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("expected 'hello', got %q", m["msg"])
	}
}

func TestApply_CollapseSpaces(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "collapse_spaces"}})
	out := s.Apply(`{"msg":"hello   world  foo"}`)
	m := decode(t, out)
	if m["msg"] != "hello world foo" {
		t.Errorf("expected 'hello world foo', got %q", m["msg"])
	}
}

func TestApply_Alphanumeric(t *testing.T) {
	s := New([]Rule{{Field: "id", Mode: "alphanumeric"}})
	out := s.Apply(`{"id":"abc-123!@#"}`)
	m := decode(t, out)
	if m["id"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", m["id"])
	}
}

func TestApply_Fallback(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "trim", Fallback: "(empty)"}})
	out := s.Apply(`{"msg":"   "}`)
	m := decode(t, out)
	if m["msg"] != "(empty)" {
		t.Errorf("expected fallback '(empty)', got %q", m["msg"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	s := New([]Rule{{Field: "missing", Mode: "trim"}})
	input := `{"msg":"hello"}`
	out := s.Apply(input)
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("expected msg unchanged, got %q", m["msg"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "trim"}})
	input := `not-json`
	if got := s.Apply(input); got != input {
		t.Errorf("expected passthrough for invalid JSON, got %q", got)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	s := New([]Rule{{Field: "count", Mode: "trim"}})
	input := `{"count":42}`
	out := s.Apply(input)
	m := decode(t, out)
	if v, ok := m["count"].(float64); !ok || v != 42 {
		t.Errorf("expected count=42 unchanged, got %v", m["count"])
	}
}
