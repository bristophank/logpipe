package spliceor

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	s := New(nil)
	line := `{"msg":"hello"}`
	out, err := s.Apply(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"src"}, Position: "after"}})
	out, err := s.Apply("not-json")
	if err != nil {
		t.Fatal(err)
	}
	if out != "not-json" {
		t.Errorf("expected passthrough for invalid JSON, got %s", out)
	}
}

func TestApply_After(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"level"}, Position: "after", Sep: "-"}})
	out, err := s.Apply(`{"msg":"hello","level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "hello-info" {
		t.Errorf("unexpected msg: %v", m["msg"])
	}
}

func TestApply_Before(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"level"}, Position: "before", Sep: ": "}})
	out, err := s.Apply(`{"msg":"hello","level":"warn"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "warn: hello" {
		t.Errorf("unexpected msg: %v", m["msg"])
	}
}

func TestApply_Replace(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"code", "text"}, Position: "replace", Sep: "|"}})
	out, err := s.Apply(`{"msg":"old","code":"404","text":"not found"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "404|not found" {
		t.Errorf("unexpected msg: %v", m["msg"])
	}
}

func TestApply_MissingSourceField_Skipped(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"missing"}, Position: "after", Sep: "-"}})
	out, err := s.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	// no missing field, injected is empty, result is "hello-"
	if m["msg"] != "hello-" {
		t.Errorf("unexpected msg: %v", m["msg"])
	}
}

func TestApply_DefaultSepIsSpace(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"level"}, Position: "after"}})
	out, err := s.Apply(`{"msg":"hi","level":"debug"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "hi debug" {
		t.Errorf("unexpected msg: %v", m["msg"])
	}
}
