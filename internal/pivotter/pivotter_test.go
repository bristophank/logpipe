package pivotter

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
	p := New(nil)
	input := `{"metrics":[{"name":"cpu","value":0.9}]}`
	if got := p.Apply(input); got != input {
		t.Fatalf("expected passthrough, got %s", got)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p := New([]Rule{{Source: "metrics", KeyField: "name", ValueField: "value"}})
	input := `not json`
	if got := p.Apply(input); got != input {
		t.Fatalf("expected passthrough on invalid JSON, got %s", got)
	}
}

func TestApply_PivotsArray(t *testing.T) {
	p := New([]Rule{
		{Source: "metrics", KeyField: "name", ValueField: "value"},
	})
	input := `{"host":"srv1","metrics":[{"name":"cpu","value":0.9},{"name":"mem","value":512}]}`
	out := decode(t, p.Apply(input))

	if out["cpu"] != 0.9 {
		t.Fatalf("expected cpu=0.9, got %v", out["cpu"])
	}
	if out["mem"] != float64(512) {
		t.Fatalf("expected mem=512, got %v", out["mem"])
	}
	if out["metrics"] == nil {
		t.Fatal("expected metrics field to remain")
	}
}

func TestApply_WithPrefix(t *testing.T) {
	p := New([]Rule{
		{Source: "tags", KeyField: "k", ValueField: "v", Prefix: "tag_"},
	})
	input := `{"tags":[{"k":"env","v":"prod"}]}`
	out := decode(t, p.Apply(input))

	if out["tag_env"] != "prod" {
		t.Fatalf("expected tag_env=prod, got %v", out["tag_env"])
	}
}

func TestApply_DeleteSource(t *testing.T) {
	p := New([]Rule{
		{Source: "metrics", KeyField: "name", ValueField: "value", DeleteSource: true},
	})
	input := `{"metrics":[{"name":"cpu","value":0.5}]}`
	out := decode(t, p.Apply(input))

	if _, exists := out["metrics"]; exists {
		t.Fatal("expected metrics field to be deleted")
	}
	if out["cpu"] != 0.5 {
		t.Fatalf("expected cpu=0.5, got %v", out["cpu"])
	}
}

func TestApply_MissingSource(t *testing.T) {
	p := New([]Rule{
		{Source: "metrics", KeyField: "name", ValueField: "value"},
	})
	input := `{"host":"srv1"}`
	out := decode(t, p.Apply(input))

	if out["host"] != "srv1" {
		t.Fatalf("expected host unchanged, got %v", out["host"])
	}
}

func TestApply_SkipsItemsMissingKeyField(t *testing.T) {
	p := New([]Rule{
		{Source: "metrics", KeyField: "name", ValueField: "value"},
	})
	input := `{"metrics":[{"value":1},{"name":"cpu","value":0.8}]}`
	out := decode(t, p.Apply(input))

	if out["cpu"] != 0.8 {
		t.Fatalf("expected cpu=0.8, got %v", out["cpu"])
	}
}
