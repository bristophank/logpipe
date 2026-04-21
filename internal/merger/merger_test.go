package merger

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	m := New(nil)
	line := `{"a":"foo","b":"bar"}`
	out, err := m.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_MergesFields(t *testing.T) {
	m := New([]Rule{
		{Fields: []string{"first", "last"}, Target: "full_name"},
	})
	out, err := m.Apply(`{"first":"John","last":"Doe"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if obj["full_name"] != "John Doe" {
		t.Errorf("expected 'John Doe', got %v", obj["full_name"])
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	m := New([]Rule{
		{Fields: []string{"host", "port"}, Target: "address", Separator: ":"},
	})
	out, err := m.Apply(`{"host":"localhost","port":"8080"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if obj["address"] != "localhost:8080" {
		t.Errorf("expected 'localhost:8080', got %v", obj["address"])
	}
}

func TestApply_DeleteSrc(t *testing.T) {
	m := New([]Rule{
		{Fields: []string{"first", "last"}, Target: "full_name", DeleteSrc: true},
	})
	out, err := m.Apply(`{"first":"Jane","last":"Smith"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if _, ok := obj["first"]; ok {
		t.Error("expected 'first' to be deleted")
	}
	if _, ok := obj["last"]; ok {
		t.Error("expected 'last' to be deleted")
	}
	if obj["full_name"] != "Jane Smith" {
		t.Errorf("expected 'Jane Smith', got %v", obj["full_name"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	m := New([]Rule{
		{Fields: []string{"a", "b"}, Target: "merged"},
	})
	out, err := m.Apply(`{"a":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj := decode(t, out)
	if obj["merged"] != "hello" {
		t.Errorf("expected 'hello', got %v", obj["merged"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m := New([]Rule{
		{Fields: []string{"a", "b"}, Target: "c"},
	})
	line := `not-json`
	out, err := m.Apply(line)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if out != line {
		t.Errorf("expected original line on error, got %s", out)
	}
}
