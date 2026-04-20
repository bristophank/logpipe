package censor

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
	c := New(nil)
	line := `{"user":"alice","action":"login"}`
	out, err := c.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_CensorsMatchingValue(t *testing.T) {
	c := New([]Rule{
		{Field: "user", Values: []string{"admin"}, Mask: "[HIDDEN]"},
	})
	out, err := c.Apply(`{"user":"admin","action":"delete"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["user"] != "[HIDDEN]" {
		t.Errorf("expected [HIDDEN], got %v", m["user"])
	}
}

func TestApply_DefaultMask(t *testing.T) {
	c := New([]Rule{
		{Field: "level", Values: []string{"secret"}},
	})
	out, err := c.Apply(`{"level":"secret"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != defaultMask {
		t.Errorf("expected %s, got %v", defaultMask, m["level"])
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	c := New([]Rule{
		{Field: "role", Values: []string{"Admin"}},
	})
	out, err := c.Apply(`{"role":"ADMIN"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["role"] != defaultMask {
		t.Errorf("expected mask, got %v", m["role"])
	}
}

func TestApply_NoMatchPassesThrough(t *testing.T) {
	c := New([]Rule{
		{Field: "user", Values: []string{"root"}},
	})
	line := `{"user":"alice"}`
	out, err := c.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["user"] != "alice" {
		t.Errorf("expected alice, got %v", decode(t, out)["user"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c := New([]Rule{
		{Field: "secret", Values: []string{"password"}},
	})
	line := `{"user":"bob"}`
	out, err := c.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decode(t, out)["user"] != "bob" {
		t.Errorf("expected bob")
	}
}

func TestApply_InvalidJSON_ReturnsError(t *testing.T) {
	c := New([]Rule{{Field: "x", Values: []string{"y"}}})
	out, err := c.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if out != "not-json" {
		t.Errorf("expected original line on error")
	}
}
