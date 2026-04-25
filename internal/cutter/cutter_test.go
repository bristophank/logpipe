package cutter

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return obj
}

func TestApply_NoRules(t *testing.T) {
	c := New(nil)
	in := `{"msg":"hello world"}`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Fatalf("expected passthrough, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c := New([]Rule{{Field: "msg", Start: 0, End: 3}})
	out, err := c.Apply("not-json")
	if err == nil {
		t.Fatal("expected error")
	}
	if out != "not-json" {
		t.Fatalf("expected original line on error, got %s", out)
	}
}

func TestApply_BasicCut(t *testing.T) {
	c := New([]Rule{{Field: "msg", Start: 0, End: 5}})
	out, err := c.Apply(`{"msg":"hello world"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %v", obj["msg"])
	}
}

func TestApply_EndZeroMeansToEnd(t *testing.T) {
	c := New([]Rule{{Field: "msg", Start: 6, End: 0}})
	out, err := c.Apply(`{"msg":"hello world"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["msg"] != "world" {
		t.Fatalf("expected 'world', got %v", obj["msg"])
	}
}

func TestApply_WritesToAlias(t *testing.T) {
	c := New([]Rule{{Field: "msg", Start: 0, End: 5, As: "short"}})
	out, err := c.Apply(`{"msg":"hello world"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["short"] != "hello" {
		t.Fatalf("expected 'hello' in 'short', got %v", obj["short"])
	}
	if obj["msg"] != "hello world" {
		t.Fatalf("original field should be unchanged, got %v", obj["msg"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c := New([]Rule{{Field: "missing", Start: 0, End: 3}})
	in := `{"msg":"hello"}`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if _, ok := obj["missing"]; ok {
		t.Fatal("missing field should not be created")
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	c := New([]Rule{{Field: "count", Start: 0, End: 1}})
	in := `{"count":42}`
	out, err := c.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if v, _ := obj["count"].(float64); v != 42 {
		t.Fatalf("expected count=42 unchanged, got %v", obj["count"])
	}
}
