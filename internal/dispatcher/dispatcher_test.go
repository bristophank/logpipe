package dispatcher

import (
	"bytes"
	"strings"
	"testing"
)

func sinks(names ...string) map[string](*bytes.Buffer) {
	m := make(map[string]*bytes.Buffer, len(names))
	for _, n := range names {
		m[n] = &bytes.Buffer{}
	}
	return m
}

func writerMap(m map[string]*bytes.Buffer) map[string]interface{ Write([]byte) (int, error) } {
	out := make(map[string]interface{ Write([]byte) (int, error) }, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func TestDispatch_NoRules_UsesDefault(t *testing.T) {
	bufs := sinks("default")
	d, err := New(Config{DefaultSink: "default"}, map[string]interface{ Write([]byte) (int, error) }{"default": bufs["default"]})
	if err != nil {
		t.Fatal(err)
	}
	sink, err := d.Dispatch(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	if sink != "default" {
		t.Fatalf("expected default, got %q", sink)
	}
	if !strings.Contains(bufs["default"].String(), "hello") {
		t.Error("expected line in default sink")
	}
}

func TestDispatch_MatchedRule(t *testing.T) {
	bufs := sinks("errors", "default")
	sinkMap := map[string]interface{ Write([]byte) (int, error) }{"errors": bufs["errors"], "default": bufs["default"]}
	d, err := New(Config{
		Rules:       []Rule{{Field: "level", Value: "error", Sink: "errors"}},
		DefaultSink: "default",
	}, sinkMap)
	if err != nil {
		t.Fatal(err)
	}
	sink, _ := d.Dispatch(`{"level":"error","msg":"boom"}`)
	if sink != "errors" {
		t.Fatalf("expected errors sink, got %q", sink)
	}
	if bufs["default"].Len() != 0 {
		t.Error("default sink should be empty")
	}
}

func TestDispatch_NoMatch_NoDefault_Drops(t *testing.T) {
	bufs := sinks("errors")
	sinkMap := map[string]interface{ Write([]byte) (int, error) }{"errors": bufs["errors"]}
	d, _ := New(Config{Rules: []Rule{{Field: "level", Value: "error", Sink: "errors"}}}, sinkMap)
	sink, err := d.Dispatch(`{"level":"info","msg":"hi"}`)
	if err != nil {
		t.Fatal(err)
	}
	if sink != "" {
		t.Fatalf("expected empty sink (drop), got %q", sink)
	}
}

func TestDispatch_EmptyLine_Ignored(t *testing.T) {
	bufs := sinks("default")
	sinkMap := map[string]interface{ Write([]byte) (int, error) }{"default": bufs["default"]}
	d, _ := New(Config{DefaultSink: "default"}, sinkMap)
	sink, err := d.Dispatch("   ")
	if err != nil {
		t.Fatal(err)
	}
	if sink != "" {
		t.Fatalf("expected empty, got %q", sink)
	}
}

func TestNew_UnknownSink_ReturnsError(t *testing.T) {
	_, err := New(
		Config{Rules: []Rule{{Field: "level", Value: "error", Sink: "missing"}}},
		map[string]interface{ Write([]byte) (int, error) }{},
	)
	if err == nil {
		t.Fatal("expected error for unknown sink")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	bufs := sinks("out")
	sinkMap := map[string]interface{ Write([]byte) (int, error) }{"out": bufs["out"]}
	_, err := New(Config{Rules: []Rule{{Field: "", Value: "x", Sink: "out"}}}, sinkMap)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
