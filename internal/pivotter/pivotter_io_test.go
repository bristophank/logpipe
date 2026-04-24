package pivotter

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_PivotsLines(t *testing.T) {
	p := New([]Rule{
		{Source: "metrics", KeyField: "name", ValueField: "value", DeleteSource: true},
	})

	input := `{"host":"a","metrics":[{"name":"cpu","value":0.5}]}
{"host":"b","metrics":[{"name":"mem","value":256}]}
`
	var buf bytes.Buffer
	if err := Stream(p, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	out0 := decode(t, lines[0])
	if out0["cpu"] != 0.5 {
		t.Fatalf("line 0: expected cpu=0.5, got %v", out0["cpu"])
	}
	if _, ok := out0["metrics"]; ok {
		t.Fatal("line 0: metrics should be deleted")
	}

	out1 := decode(t, lines[1])
	if out1["mem"] != float64(256) {
		t.Fatalf("line 1: expected mem=256, got %v", out1["mem"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	p := New([]Rule{
		{Source: "m", KeyField: "k", ValueField: "v"},
	})

	input := `{"m":[{"k":"x","v":1}]}

{"m":[{"k":"y","v":2}]}
`
	var buf bytes.Buffer
	if err := Stream(p, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 non-empty lines, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	p := New([]Rule{
		{Source: "metrics", KeyField: "name", ValueField: "value"},
	})

	input := "not json\n"
	var buf bytes.Buffer
	if err := Stream(p, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	if got := strings.TrimSpace(buf.String()); got != "not json" {
		t.Fatalf("expected passthrough of invalid JSON, got %q", got)
	}
}
