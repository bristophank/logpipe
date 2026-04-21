package extractor

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
	e := New(nil)
	out, err := e.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `{"level":"info"}` {
		t.Errorf("expected unchanged line, got %s", out)
	}
}

func TestApply_BasicExtraction(t *testing.T) {
	e := New([]Rule{{Field: "msg", Target: "msg_copy"}})
	out, err := e.Apply(`{"msg":"hello world"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["msg_copy"] != "hello world" {
		t.Errorf("expected msg_copy=hello world, got %v", m["msg_copy"])
	}
}

func TestApply_PrefixTrim(t *testing.T) {
	e := New([]Rule{{Field: "path", Target: "short_path", Prefix: "/api/v1/"}})
	out, err := e.Apply(`{"path":"/api/v1/users"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["short_path"] != "users" {
		t.Errorf("expected short_path=users, got %v", m["short_path"])
	}
}

func TestApply_SuffixTrim(t *testing.T) {
	e := New([]Rule{{Field: "filename", Target: "base", Suffix: ".log"}})
	out, err := e.Apply(`{"filename":"app.log"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["base"] != "app" {
		t.Errorf("expected base=app, got %v", m["base"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	e := New([]Rule{{Field: "nonexistent", Target: "out"}})
	out, err := e.Apply(`{"level":"warn"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["out"]; ok {
		t.Error("expected 'out' field to be absent")
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	e := New([]Rule{{Field: "msg", Target: "msg_copy"}})
	input := `not json`
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_DefaultTarget(t *testing.T) {
	e := New([]Rule{{Field: "status"}})
	out, err := e.Apply(`{"status":"200 OK"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["status_extracted"] != "200 OK" {
		t.Errorf("expected status_extracted=200 OK, got %v", m["status_extracted"])
	}
}
