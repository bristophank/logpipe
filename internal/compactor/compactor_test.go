package compactor

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
	line := `{"a":"hello","b":""}`
	out, err := c.Apply(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	c := New([]Rule{{DropEmpty: true}})
	out, err := c.Apply("not-json")
	if err != nil {
		t.Fatal(err)
	}
	if out != "not-json" {
		t.Errorf("expected passthrough for invalid JSON, got %s", out)
	}
}

func TestApply_DropEmptyString_AllFields(t *testing.T) {
	c := New([]Rule{{DropEmpty: true}})
	out, err := c.Apply(`{"a":"hello","b":"","c":"world"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["b"]; ok {
		t.Error("expected 'b' to be removed")
	}
	if m["a"] != "hello" || m["c"] != "world" {
		t.Error("non-empty fields should remain")
	}
}

func TestApply_DropNull_AllFields(t *testing.T) {
	c := New([]Rule{{DropNull: true}})
	out, err := c.Apply(`{"a":null,"b":"keep"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["a"]; ok {
		t.Error("expected 'a' to be removed")
	}
	if m["b"] != "keep" {
		t.Error("'b' should remain")
	}
}

func TestApply_DropZero_AllFields(t *testing.T) {
	c := New([]Rule{{DropZero: true}})
	out, err := c.Apply(`{"count":0,"score":3.5}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["count"]; ok {
		t.Error("expected 'count' to be removed")
	}
	if m["score"] != 3.5 {
		t.Error("'score' should remain")
	}
}

func TestApply_DropFalse_SpecificField(t *testing.T) {
	c := New([]Rule{{Field: "active", DropFalse: true}})
	out, err := c.Apply(`{"active":false,"ready":false}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["active"]; ok {
		t.Error("expected 'active' to be removed")
	}
	if _, ok := m["ready"]; !ok {
		t.Error("'ready' should remain (not targeted)")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	c := New([]Rule{
		{DropNull: true},
		{DropEmpty: true},
	})
	out, err := c.Apply(`{"a":null,"b":"","c":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["a"]; ok {
		t.Error("'a' should be removed")
	}
	if _, ok := m["b"]; ok {
		t.Error("'b' should be removed")
	}
	if m["c"] != "ok" {
		t.Error("'c' should remain")
	}
}
