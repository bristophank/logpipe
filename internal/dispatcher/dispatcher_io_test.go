package dispatcher

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_DispatchesLines(t *testing.T) {
	errBuf := &bytes.Buffer{}
	infoBuf := &bytes.Buffer{}
	sinkMap := map[string]interface{ Write([]byte) (int, error) }{
		"errors": errBuf,
		"info":   infoBuf,
	}
	d, err := New(Config{
		Rules: []Rule{
			{Field: "level", Value: "error", Sink: "errors"},
			{Field: "level", Value: "info", Sink: "info"},
		},
	}, sinkMap)
	if err != nil {
		t.Fatal(err)
	}

	input := strings.Join([]string{
		`{"level":"error","msg":"bad"}`,
		`{"level":"info","msg":"ok"}`,
		`{"level":"error","msg":"also bad"}`,
	}, "\n")

	n, err := Stream(strings.NewReader(input), d)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3 dispatched, got %d", n)
	}
	if !strings.Contains(errBuf.String(), "bad") {
		t.Error("error sink missing expected content")
	}
	if !strings.Contains(infoBuf.String(), "ok") {
		t.Error("info sink missing expected content")
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	defBuf := &bytes.Buffer{}
	sinkMap := map[string]interface{ Write([]byte) (int, error) }{"default": defBuf}
	d, _ := New(Config{DefaultSink: "default"}, sinkMap)

	input := "\n\n" + `{"msg":"hello"}` + "\n\n"
	n, err := Stream(strings.NewReader(input), d)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected 1 dispatched, got %d", n)
	}
}
