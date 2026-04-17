package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cfg*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	p := writeTemp(t, `{"sinks":[{"name":"out","type":"stdout"}],"routes":[{"sink":"out"}],"transformers":[{"field":"level","op":"uppercase"}]}`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sinks) != 1 || cfg.Sinks[0].Name != "out" {
		t.Error("sinks not loaded")
	}
	if len(cfg.Transformers) != 1 || cfg.Transformers[0].Op != "uppercase" {
		t.Error("transformers not loaded")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := writeTemp(t, `{bad json}`)
	_, err := Load(p)
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestLoad_UnknownSinkInRoute(t *testing.T) {
	p := writeTemp(t, `{"sinks":[{"name":"out","type":"stdout"}],"routes":[{"sink":"missing"}]}`)
	_, err := Load(p)
	if err == nil {
		t.Error("expected error for unknown sink")
	}
}

func TestLoad_UnknownTransformerOp(t *testing.T) {
	p := writeTemp(t, `{"sinks":[],"routes":[],"transformers":[{"field":"x","op":"explode"}]}`)
	_, err := Load(p)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestLoad_EmptyConfig(t *testing.T) {
	p := writeTemp(t, `{}`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sinks) != 0 || len(cfg.Routes) != 0 {
		t.Error("expected empty config")
	}
}
