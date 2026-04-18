package parser

import (
	"testing"
)

func TestParse_JSON(t *testing.T) {
	p := New(FormatJSON)
	m, err := p.Parse(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	p := New(FormatJSON)
	_, err := p.Parse(`not json`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_Logfmt(t *testing.T) {
	p := New(FormatLogfmt)
	m, err := p.Parse(`level=info msg="hello world"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" {
		t.Fatalf("expected level=info, got %v", m["level"])
	}
}

func TestParse_LogfmtFlag(t *testing.T) {
	p := New(FormatLogfmt)
	m, err := p.Parse(`debug level=warn`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["debug"] != true {
		t.Fatalf("expected debug flag")
	}
}

func TestParse_Auto_JSON(t *testing.T) {
	p := New(FormatAuto)
	m, err := p.Parse(`{"x":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["x"]; !ok {
		t.Fatal("expected key x")
	}
}

func TestParse_Auto_Logfmt(t *testing.T) {
	p := New(FormatAuto)
	m, err := p.Parse(`a=1 b=2`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["a"] != "1" || m["b"] != "2" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParse_EmptyLine(t *testing.T) {
	p := New(FormatAuto)
	_, err := p.Parse("   ")
	if err == nil {
		t.Fatal("expected error for empty line")
	}
}

func TestNew_DefaultFormat(t *testing.T) {
	p := New("")
	if p.format != FormatAuto {
		t.Fatalf("expected auto, got %v", p.format)
	}
}
