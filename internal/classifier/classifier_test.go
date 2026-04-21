package classifier

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
	c := New("", nil)
	in := `{"level":"error"}`
	out, err := c.Apply(in)
	if err != nil || out != in {
		t.Fatalf("expected passthrough, got %q %v", out, err)
	}
}

func TestApply_EqualsMatch(t *testing.T) {
	c := New("category", []Rule{
		{Field: "level", Equals: "error", Category: "critical"},
	})
	out, err := c.Apply(`{"level":"error","msg":"oops"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["category"] != "critical" {
		t.Fatalf("expected critical, got %v", m["category"])
	}
}

func TestApply_ContainsMatch(t *testing.T) {
	c := New("", []Rule{
		{Field: "msg", Contains: "timeout", Category: "network"},
	})
	out, err := c.Apply(`{"msg":"connection timeout reached"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["category"] != "network" {
		t.Fatalf("expected network, got %v", m["category"])
	}
}

func TestApply_NoMatch_Unchanged(t *testing.T) {
	c := New("", []Rule{
		{Field: "level", Equals: "debug", Category: "verbose"},
	})
	in := `{"level":"info"}`
	out, err := c.Apply(in)
	if err != nil || out != in {
		t.Fatalf("expected unchanged, got %q %v", out, err)
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c := New("", []Rule{
		{Field: "status", Equals: "500", Category: "server_error"},
	})
	in := `{"msg":"hello"}`
	out, err := c.Apply(in)
	if err != nil || out != in {
		t.Fatalf("expected unchanged, got %q %v", out, err)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c := New("", []Rule{{Field: "x", Equals: "y", Category: "z"}})
	in := `not json`
	out, err := c.Apply(in)
	if err == nil {
		t.Fatal("expected error")
	}
	if out != in {
		t.Fatalf("expected original line returned, got %q", out)
	}
}

func TestApply_CustomOutKey(t *testing.T) {
	c := New("_class", []Rule{
		{Field: "level", Equals: "warn", Category: "warning"},
	})
	out, err := c.Apply(`{"level":"warn"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["_class"] != "warning" {
		t.Fatalf("expected warning in _class, got %v", m["_class"])
	}
}
