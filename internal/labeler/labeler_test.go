package labeler_test

import (
	"encoding/json"
	"testing"

	"github.com/logpipe/logpipe/internal/labeler"
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
	l := labeler.New(nil)
	line := `{"msg":"hello"}`
	got, err := l.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestApply_AddsLabels(t *testing.T) {
	l := labeler.New([]labeler.Rule{
		{Key: "env", Value: "production"},
		{Key: "service", Value: "api"},
	})
	got, err := l.Apply(`{"msg":"started"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, got)
	if m["env"] != "production" {
		t.Errorf("expected env=production, got %v", m["env"])
	}
	if m["service"] != "api" {
		t.Errorf("expected service=api, got %v", m["service"])
	}
	if m["msg"] != "started" {
		t.Errorf("original field lost")
	}
}

func TestApply_OverwritesExistingKey(t *testing.T) {
	l := labeler.New([]labeler.Rule{{Key: "env", Value: "staging"}})
	got, err := l.Apply(`{"env":"dev","msg":"test"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, got)
	if m["env"] != "staging" {
		t.Errorf("expected env overwritten to staging, got %v", m["env"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	l := labeler.New([]labeler.Rule{{Key: "env", Value: "prod"}})
	orig := `not-json`
	got, err := l.Apply(orig)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if got != orig {
		t.Errorf("expected original line on error, got %q", got)
	}
}

func TestApply_EmptyLine(t *testing.T) {
	l := labeler.New([]labeler.Rule{{Key: "env", Value: "prod"}})
	got, err := l.Apply("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestNew_SkipsEmptyKeys(t *testing.T) {
	l := labeler.New([]labeler.Rule{
		{Key: "", Value: "ignored"},
		{Key: "region", Value: "us-east-1"},
	})
	got, err := l.Apply(`{"msg":"ok"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, got)
	if _, ok := m[""]; ok {
		t.Error("empty key should not be present")
	}
	if m["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", m["region"])
	}
}
