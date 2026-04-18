package parser

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStream_JSON(t *testing.T) {
	input := `{"level":"info"}
{"level":"warn"}
`
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestStream_SkipsInvalid(t *testing.T) {
	input := "not json\n{\"ok\":true}\n"
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]any
	if err := json.NewDecoder(&buf).Decode(&m); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if m["ok"] != true {
		t.Fatalf("expected ok=true")
	}
}

func TestStream_Logfmt(t *testing.T) {
	input := "level=info msg=started\nlevel=error msg=failed\n"
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, FormatLogfmt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}
