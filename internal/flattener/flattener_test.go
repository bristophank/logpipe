package flattener

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
	f := New(nil)
	input := `{"a":{"b":1}}`
	if got := f.Apply(input); got != input {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestApply_FlattenAllFields(t *testing.T) {
	f := New([]Rule{{Separator: "."}})
	input := `{"a":{"b":1,"c":{"d":2}},"e":3}`
	m := decode(t, f.Apply(input))
	if m["a.b"] != float64(1) {
		t.Errorf("expected a.b=1, got %v", m["a.b"])
	}
	if m["a.c.d"] != float64(2) {
		t.Errorf("expected a.c.d=2, got %v", m["a.c.d"])
	}
	if m["e"] != float64(3) {
		t.Errorf("expected e=3, got %v", m["e"])
	}
	if _, ok := m["a"]; ok {
		t.Error("expected 'a' to be removed")
	}
}

func TestApply_FlattenSpecificField(t *testing.T) {
	f := New([]Rule{{Separator: "_", Fields: []string{"meta"}}})
	input := `{"level":"info","meta":{"host":"srv1","region":"us-east"}}`
	m := decode(t, f.Apply(input))
	if m["meta_host"] != "srv1" {
		t.Errorf("expected meta_host=srv1, got %v", m["meta_host"])
	}
	if m["meta_region"] != "us-east" {
		t.Errorf("expected meta_region=us-east, got %v", m["meta_region"])
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
	if _, ok := m["meta"]; ok {
		t.Error("expected 'meta' to be removed")
	}
}

func TestApply_WithPrefix(t *testing.T) {
	f := New([]Rule{{Prefix: "log", Separator: "."}})
	input := `{"a":{"b":"x"}}`
	m := decode(t, f.Apply(input))
	if m["log.a.b"] != "x" {
		t.Errorf("expected log.a.b=x, got %v", m["log.a.b"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	f := New([]Rule{{Separator: "."}})
	input := `not json`
	if got := f.Apply(input); got != input {
		t.Errorf("expected passthrough for invalid JSON")
	}
}

func TestApply_EmptyLine(t *testing.T) {
	f := New([]Rule{{Separator: "."}})
	if got := f.Apply(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestApply_DefaultSeparator(t *testing.T) {
	f := New([]Rule{{}})
	input := `{"x":{"y":true}}`
	m := decode(t, f.Apply(input))
	if m["x.y"] != true {
		t.Errorf("expected x.y=true with default separator, got %v", m["x.y"])
	}
}
