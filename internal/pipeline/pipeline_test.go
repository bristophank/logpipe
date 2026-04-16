package pipeline

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/router"
)

func newRouter(t *testing.T) (*router.Router, *bytes.Buffer) {
	t.Helper()
	r := router.New()
	var buf bytes.Buffer
	r.AddSink("out", &buf)
	return r, &buf
}

func TestPipeline_PassesMatchingLines(t *testing.T) {
	r, buf := newRouter(t)
	p, err := New(Config{
		SinkName: "out",
		Rules:    []filter.Rule{{Field: "level", Op: "eq", Value: "error"}},
	}, r)
	if err != nil {
		t.Fatal(err)
	}

	src := strings.NewReader(`{"level":"error","msg":"oops"}
{"level":"info","msg":"ok"}
`)
	if err := p.Run(src); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "oops") {
		t.Error("expected error line to be routed")
	}
	if strings.Contains(buf.String(), "ok") {
		t.Error("expected info line to be filtered out")
	}
}

func TestPipeline_NoRulesPassesAll(t *testing.T) {
	r, buf := newRouter(t)
	p, err := New(Config{SinkName: "out", Rules: nil}, r)
	if err != nil {
		t.Fatal(err)
	}

	src := strings.NewReader(`{"level":"debug"}
{"level":"info"}
`)
	if err := p.Run(src); err != nil {
		t.Fatal(err)
	}

	if strings.Count(buf.String(), "\n") != 2 {
		t.Errorf("expected 2 lines, got: %q", buf.String())
	}
}
