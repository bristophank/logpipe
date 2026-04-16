package formatter

import (
	"strings"
	"testing"
)

func TestNew_DefaultsToJSON(t *testing.T) {
	f := New("unknown")
	if f.format != FormatJSON {
		t.Fatalf("expected json, got %s", f.format)
	}
}

func TestFormat_JSON_Passthrough(t *testing.T) {
	f := New("json")
	input := `{"level":"info","msg":"hello"}`
	out, err := f.Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if out != input {
		t.Fatalf("expected passthrough, got %s", out)
	}
}

func TestFormat_Pretty(t *testing.T) {
	f := New("pretty")
	input := `{"level":"info","msg":"hello"}`
	out, err := f.Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "\n") {
		t.Fatal("expected indented output")
	}
	if !strings.Contains(out, "level") {
		t.Fatal("expected key 'level' in output")
	}
}

func TestFormat_Text(t *testing.T) {
	f := New("text")
	input := `{"level":"info","msg":"hello"}`
	out, err := f.Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "level=info") {
		t.Fatalf("expected level=info in %q", out)
	}
	if !strings.Contains(out, "msg=hello") {
		t.Fatalf("expected msg=hello in %q", out)
	}
}

func TestFormat_Text_KeysSorted(t *testing.T) {
	f := New("text")
	input := `{"z":"last","a":"first"}`
	out, err := f.Format(input)
	if err != nil {
		t.Fatal(err)
	}
	idxA := strings.Index(out, "a=")
	idxZ := strings.Index(out, "z=")
	if idxA > idxZ {
		t.Fatalf("expected 'a' before 'z' in %q", out)
	}
}

func TestFormat_InvalidJSON(t *testing.T) {
	for _, fmt := range []string{"pretty", "text"} {
		f := New(fmt)
		_, err := f.Format("not-json")
		if err == nil {
			t.Fatalf("expected error for format %s with invalid json", fmt)
		}
	}
}
