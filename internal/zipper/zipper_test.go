package zipper_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/zipper"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	z := zipper.New(nil)
	line := `{"a":1,"b":2}`
	out, err := z.Apply(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_ZipsFields(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"first", "last"}, Target: "name", Sep: " "},
	})
	out, err := z.Apply(`{"first":"John","last":"Doe"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["name"] != "John Doe" {
		t.Errorf("expected 'John Doe', got %v", m["name"])
	}
}

func TestApply_MissingKey_Skipped(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"a", "b"}, Target: "ab", Sep: "-"},
	})
	out, err := z.Apply(`{"a":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["ab"] != "hello" {
		t.Errorf("expected 'hello', got %v", m["ab"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"a", "b"}, Target: "ab", Sep: "-"},
	})
	_, err := z.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_DefaultSep(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"x", "y"}, Target: "xy"},
	})
	out, err := z.Apply(`{"x":"foo","y":"bar"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if !strings.Contains(m["xy"].(string), "foo") {
		t.Errorf("unexpected value: %v", m["xy"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"a", "b"}, Target: "ab", Sep: "_"},
		{Keys: []string{"c", "d"}, Target: "cd", Sep: "."},
	})
	out, err := z.Apply(`{"a":"1","b":"2","c":"x","d":"y"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["ab"] != "1_2" {
		t.Errorf("expected '1_2', got %v", m["ab"])
	}
	if m["cd"] != "x.y" {
		t.Errorf("expected 'x.y', got %v", m["cd"])
	}
}
