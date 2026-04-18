package tagger

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
	tg := New(nil)
	out, err := tg.Apply(`{"level":"error"}`)
	if err != nil || out != `{"level":"error"}` {
		t.Fatalf("unexpected: %v %v", out, err)
	}
}

func TestApply_MatchingRule(t *testing.T) {
	tg := New([]Rule{{Field: "level", Value: "error", Tag: "alert", TagValue: "true"}})
	out, err := tg.Apply(`{"level":"error","msg":"boom"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["alert"] != "true" {
		t.Fatalf("expected alert=true, got %v", m["alert"])
	}
}

func TestApply_NonMatchingRule(t *testing.T) {
	tg := New([]Rule{{Field: "level", Value: "error", Tag: "alert", TagValue: "true"}})
	out, err := tg.Apply(`{"level":"info","msg":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["alert"]; ok {
		t.Fatal("alert should not be set")
	}
}

func TestApply_MissingField(t *testing.T) {
	tg := New([]Rule{{Field: "level", Value: "error", Tag: "alert", TagValue: "true"}})
	out, err := tg.Apply(`{"msg":"no level"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["alert"]; ok {
		t.Fatal("alert should not be set")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tg := New([]Rule{{Field: "level", Value: "error", Tag: "alert", TagValue: "true"}})
	out, err := tg.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error")
	}
	if out != `not-json` {
		t.Fatalf("expected original line back, got %v", out)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	tg := New([]Rule{
		{Field: "level", Value: "error", Tag: "severity", TagValue: "high"},
		{Field: "env", Value: "prod", Tag: "critical", TagValue: "yes"},
	})
	out, err := tg.Apply(`{"level":"error","env":"prod"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["severity"] != "high" {
		t.Fatalf("expected severity=high")
	}
	if m["critical"] != "yes" {
		t.Fatalf("expected critical=yes")
	}
}
