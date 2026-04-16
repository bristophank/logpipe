package config

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cfg*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `{
		"sinks": [{"name": "out", "path": "/tmp/out.log"}],
		"routes": [{"sink": "out", "rules": [{"field": "level", "operator": "eq", "value": "error"}]}]
	}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sinks) != 1 || cfg.Sinks[0].Name != "out" {
		t.Errorf("unexpected sinks: %+v", cfg.Sinks)
	}
	if len(cfg.Routes) != 1 || cfg.Routes[0].Sink != "out" {
		t.Errorf("unexpected routes: %+v", cfg.Routes)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/cfg.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `not json`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoad_UnknownSinkInRoute(t *testing.T) {
	path := writeTemp(t, `{
		"sinks": [{"name": "out", "path": "/tmp/out.log"}],
		"routes": [{"sink": "missing", "rules": []}]
	}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for unknown sink")
	}
}

func TestLoad_SinkMissingName(t *testing.T) {
	path := writeTemp(t, `{
		"sinks": [{"path": "/tmp/out.log"}],
		"routes": []
	}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for sink missing name")
	}
}
