package config

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "logpipe-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `{"format":"pretty","sample_rate":3,"sinks":[{"name":"out","target":"stdout"}],"routes":[{"sink":"out","rules":[]}]}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "pretty" {
		t.Errorf("expected format pretty, got %s", cfg.Format)
	}
	if cfg.SampleRate != 3 {
		t.Errorf("expected sample_rate 3, got %d", cfg.SampleRate)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/cfg.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{bad json`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoad_UnknownSinkInRoute(t *testing.T) {
	path := writeTemp(t, `{"sinks":[{"name":"out","target":"stdout"}],"routes":[{"sink":"missing"}]}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unknown sink in route")
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTemp(t, `{"sinks":[{"name":"out","target":"stdout"}],"routes":[{"sink":"out"}]}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected default format json, got %s", cfg.Format)
	}
	if cfg.SampleRate != 1 {
		t.Errorf("expected default sample_rate 1, got %d", cfg.SampleRate)
	}
}
