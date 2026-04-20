package replay

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempReplay(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "replay-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestStream_ReadsFile(t *testing.T) {
	path := writeTempReplay(t, "{\"msg\":\"hello\"}\n{\"msg\":\"world\"}\n")
	var out bytes.Buffer
	if err := Stream(Config{}, path, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestStream_MissingFile(t *testing.T) {
	var out bytes.Buffer
	if err := Stream(Config{}, "/no/such/file.log", &out); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestStream_EmptyFile(t *testing.T) {
	path := writeTempReplay(t, "")
	var out bytes.Buffer
	if err := Stream(Config{}, path, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty output, got %q", out.String())
	}
}
