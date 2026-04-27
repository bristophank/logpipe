package expander

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
	e := New(nil)
	in := `{"tags":"a,b,c"}`
	out, err := e.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Fatalf("expected passthrough, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	e := New([]Rule{{Field: "tags", Delimiter: ","}})
	out, err := e.Apply("not-json")
	if err == nil {
		t.Fatal("expected error")
	}
	if out != "not-json" {
		t.Fatalf("expected original line on error, got %s", out)
	}
}

func TestApply_SplitsOnComma(t *testing.T) {
	e := New([]Rule{{Field: "tags", Delimiter: ","}})
	out, err := e.Apply(`{"tags":"a,b,c"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	vals, ok := m["tags"].([]any)
	if !ok {
		t.Fatalf("expected array, got %T", m["tags"])
	}
	if len(vals) != 3 || vals[0] != "a" || vals[1] != "b" || vals[2] != "c" {
		t.Fatalf("unexpected values: %v", vals)
	}
}

func TestApply_TargetField(t *testing.T) {
	e := New([]Rule{{Field: "csv", Delimiter: ";", Target: "items"}})
	out, err := e.Apply(`{"csv":"x;y"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["csv"]; !ok {
		t.Fatal("original field should remain")
	}
	vals, ok := m["items"].([]any)
	if !ok || len(vals) != 2 {
		t.Fatalf("expected 2-element array in target, got %v", m["items"])
	}
}

func TestApply_DefaultDelimiter(t *testing.T) {
	e := New([]Rule{{Field: "tags"}})
	out, err := e.Apply(`{"tags":"go,rust,zig"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	vals := m["tags"].([]any)
	if len(vals) != 3 {
		t.Fatalf("expected 3 items, got %d", len(vals))
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	e := New([]Rule{{Field: "missing", Delimiter: ","}})
	in := `{"level":"info"}`
	out, err := e.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["missing"]; ok {
		t.Fatal("missing field should not be created")
	}
}

func TestApply_TrimsWhitespace(t *testing.T) {
	e := New([]Rule{{Field: "tags", Delimiter: ","}})
	out, err := e.Apply(`{"tags":" a , b , c "}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	vals := m["tags"].([]any)
	if vals[0] != "a" || vals[1] != "b" || vals[2] != "c" {
		t.Fatalf("expected trimmed values, got %v", vals)
	}
}
