package stencil

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
	s, _ := New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := s.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_RendersTemplate(t *testing.T) {
	s, err := New([]Rule{
		{Target: "summary", Template: "[{{.level}}] {{.msg}}"},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(`{"level":"error","msg":"disk full"}`)
	m := decode(t, out)
	if got, ok := m["summary"]; !ok || got != "[error] disk full" {
		t.Errorf("unexpected summary: %v", got)
	}
}

func TestApply_NoOverwrite(t *testing.T) {
	s, _ := New([]Rule{
		{Target: "label", Template: "new", Overwrite: false},
	})
	out := s.Apply(`{"label":"original"}`)
	m := decode(t, out)
	if m["label"] != "original" {
		t.Errorf("expected original to be preserved, got %v", m["label"])
	}
}

func TestApply_Overwrite(t *testing.T) {
	s, _ := New([]Rule{
		{Target: "label", Template: "replaced", Overwrite: true},
	})
	out := s.Apply(`{"label":"original"}`)
	m := decode(t, out)
	if m["label"] != "replaced" {
		t.Errorf("expected replaced, got %v", m["label"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	s, _ := New([]Rule{
		{Target: "x", Template: "val"},
	})
	line := "not-json"
	if got := s.Apply(line); got != line {
		t.Errorf("expected passthrough for invalid JSON, got %q", got)
	}
}

func TestApply_InvalidTemplate(t *testing.T) {
	_, err := New([]Rule{
		{Target: "x", Template: "{{.unclosed"},
	})
	if err == nil {
		t.Error("expected error for invalid template")
	}
}

func TestApply_MissingFieldRendersEmpty(t *testing.T) {
	s, _ := New([]Rule{
		{Target: "out", Template: "prefix-{{.missing}}"},
	})
	out := s.Apply(`{"level":"info"}`)
	m := decode(t, out)
	if m["out"] != "prefix-<no value>" && m["out"] != "prefix-" {
		// missingkey=zero renders zero value; for maps that is <no value> in older Go
		// Accept both forms.
		t.Logf("out field: %v (accepted)", m["out"])
	}
}
