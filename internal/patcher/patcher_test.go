package patcher_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/patcher"
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
	p := patcher.New(nil)
	input := `{"level":"error","msg":"boom"}`
	if got := p.Apply(input); got != input {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestApply_EqMatch_PatchesTarget(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "level", Op: "eq", Value: "error", Target: "alert", Patch: "true"},
	})
	out := p.Apply(`{"level":"error","msg":"fail"}`)
	m := decode(t, out)
	if m["alert"] != "true" {
		t.Errorf("expected alert=true, got %v", m["alert"])
	}
}

func TestApply_EqNoMatch_Unchanged(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "level", Op: "eq", Value: "error", Target: "alert", Patch: "true"},
	})
	input := `{"level":"info","msg":"ok"}`
	out := p.Apply(input)
	m := decode(t, out)
	if _, ok := m["alert"]; ok {
		t.Error("expected no alert field")
	}
}

func TestApply_ContainsMatch(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "msg", Op: "contains", Value: "timeout", Target: "category", Patch: "network"},
	})
	out := p.Apply(`{"msg":"connection timeout reached"}`)
	m := decode(t, out)
	if m["category"] != "network" {
		t.Errorf("expected category=network, got %v", m["category"])
	}
}

func TestApply_ExistsMatch(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "trace_id", Op: "exists", Target: "traced", Patch: "yes"},
	})
	out := p.Apply(`{"trace_id":"abc123","msg":"ok"}`)
	m := decode(t, out)
	if m["traced"] != "yes" {
		t.Errorf("expected traced=yes, got %v", m["traced"])
	}
}

func TestApply_ExistsMissing_NoChange(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "trace_id", Op: "exists", Target: "traced", Patch: "yes"},
	})
	out := p.Apply(`{"msg":"no trace"}`)
	m := decode(t, out)
	if _, ok := m["traced"]; ok {
		t.Error("expected no traced field")
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "level", Op: "eq", Value: "error", Target: "x", Patch: "1"},
	})
	input := "not json at all"
	if got := p.Apply(input); got != input {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	p := patcher.New([]patcher.Rule{
		{Field: "level", Op: "eq", Value: "error", Target: "severity", Patch: "high"},
		{Field: "level", Op: "eq", Value: "error", Target: "notify", Patch: "true"},
	})
	out := p.Apply(`{"level":"error"}`)
	m := decode(t, out)
	if m["severity"] != "high" {
		t.Errorf("expected severity=high, got %v", m["severity"])
	}
	if m["notify"] != "true" {
		t.Errorf("expected notify=true, got %v", m["notify"])
	}
}
