package splitter_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/logpipe/internal/splitter"
)

func TestWrite_NoRules_UsesFallback(t *testing.T) {
	fallback := &bytes.Buffer{}
	s := splitter.New(nil, nil, fallback)
	if err := s.Write([]byte(`{"level":"info","msg":"hello"}`)); err != nil {
		t.Fatal(err)
	}
	if fallback.Len() == 0 {
		t.Fatal("expected fallback to receive line")
	}
}

func TestWrite_MatchedSink(t *testing.T) {
	errorSink := &bytes.Buffer{}
	sinks := map[string]interface{ Write([]byte) (int, error) }{
		"errors": errorSink,
	}
	rules := []splitter.Rule{
		{Field: "level", Value: "error", Sink: "errors"},
	}
	s := splitter.New(rules, map[string]interface{ Write([]byte) (int, error) }(sinks), nil)
	_ = s.Write([]byte(`{"level":"error","msg":"boom"}`))
	if errorSink.Len() == 0 {
		t.Fatal("expected errors sink to receive line")
	}
}

func TestWrite_NoMatch_NoFallback(t *testing.T) {
	rules := []splitter.Rule{
		{Field: "level", Value: "error", Sink: "errors"},
	}
	s := splitter.New(rules, map[string]interface{ Write([]byte) (int, error) }{}, nil)
	if err := s.Write([]byte(`{"level":"info","msg":"ok"}`)); err != nil {
		t.Fatal(err)
	}
}

func TestWrite_InvalidJSON(t *testing.T) {
	s := splitter.New(nil, nil, &bytes.Buffer{})
	if err := s.Write([]byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestWrite_MultipleRulesMatch(t *testing.T) {
	a, b := &bytes.Buffer{}, &bytes.Buffer{}
	sinks := map[string]interface{ Write([]byte) (int, error) }{
		"a": a, "b": b,
	}
	rules := []splitter.Rule{
		{Field: "level", Value: "warn", Sink: "a"},
		{Field: "level", Value: "warn", Sink: "b"},
	}
	s := splitter.New(rules, sinks, nil)
	_ = s.Write([]byte(`{"level":"warn","msg":"watch out"}`))
	if a.Len() == 0 || b.Len() == 0 {
		t.Fatal("expected both sinks to receive line")
	}
}

// TestWrite_FallbackNotWrittenOnMatch verifies that the fallback sink does not
// receive a line when a rule matches and routes it to a named sink.
func TestWrite_FallbackNotWrittenOnMatch(t *testing.T) {
	errorSink := &bytes.Buffer{}
	fallback := &bytes.Buffer{}
	sinks := map[string]interface{ Write([]byte) (int, error) }{
		"errors": errorSink,
	}
	rules := []splitter.Rule{
		{Field: "level", Value: "error", Sink: "errors"},
	}
	s := splitter.New(rules, sinks, fallback)
	_ = s.Write([]byte(`{"level":"error","msg":"boom"}`))
	if fallback.Len() != 0 {
		t.Fatal("expected fallback to not receive line when a rule matched")
	}
	if errorSink.Len() == 0 {
		t.Fatal("expected errors sink to receive line")
	}
}
