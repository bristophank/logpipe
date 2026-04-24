package sorter

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	s := New(nil)
	input := `{"items":[{"name":"b"},{"name":"a"}]}`
	if got := s.Apply(input); got != input {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s := New([]Rule{{Field: "items", By: "name"}})
	input := `not-json`
	if got := s.Apply(input); got != input {
		t.Errorf("expected passthrough on invalid JSON, got %s", got)
	}
}

func TestApply_SortsAscending(t *testing.T) {
	s := New([]Rule{{Field: "items", By: "name", Order: "asc"}})
	input := `{"items":[{"name":"charlie"},{"name":"alice"},{"name":"bob"}]}`
	out := decode(t, s.Apply(input))
	arr := out["items"].([]interface{})
	names := []string{
		arr[0].(map[string]interface{})["name"].(string),
		arr[1].(map[string]interface{})["name"].(string),
		arr[2].(map[string]interface{})["name"].(string),
	}
	if names[0] != "alice" || names[1] != "bob" || names[2] != "charlie" {
		t.Errorf("unexpected order: %v", names)
	}
}

func TestApply_SortsDescending(t *testing.T) {
	s := New([]Rule{{Field: "items", By: "name", Order: "desc"}})
	input := `{"items":[{"name":"alice"},{"name":"charlie"},{"name":"bob"}]}`
	out := decode(t, s.Apply(input))
	arr := out["items"].([]interface{})
	first := arr[0].(map[string]interface{})["name"].(string)
	if first != "charlie" {
		t.Errorf("expected charlie first in desc order, got %s", first)
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	s := New([]Rule{{Field: "missing", By: "name"}})
	input := `{"level":"info"}`
	out := s.Apply(input)
	m := decode(t, out)
	if _, ok := m["missing"]; ok {
		t.Error("unexpected field inserted")
	}
}

func TestApply_FieldNotArray_Unchanged(t *testing.T) {
	s := New([]Rule{{Field: "level", By: "name"}})
	input := `{"level":"info"}`
	out := decode(t, s.Apply(input))
	if out["level"] != "info" {
		t.Errorf("expected level=info unchanged, got %v", out["level"])
	}
}

func TestApply_EmptyLine_Passthrough(t *testing.T) {
	s := New([]Rule{{Field: "items", By: "name"}})
	if got := s.Apply(""); got != "" {
		t.Errorf("expected empty passthrough, got %q", got)
	}
}
