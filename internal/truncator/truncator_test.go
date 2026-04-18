package truncator

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	tr := New(nil)
	line := `{"msg":"hello world"}`
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestApply_ShortField_Unchanged(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 20}})
	line := `{"msg":"hi"}`
	if got := tr.Apply(line); decode(t, got)["msg"] != "hi" {
		t.Errorf("unexpected truncation")
	}
}

func TestApply_LongField_Truncated(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 5}})
	line := `{"msg":"hello world"}`
	got := decode(t, tr.Apply(line))
	if got["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", got["msg"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	tr := New([]Rule{{Field: "missing", MaxLen: 3}})
	line := `{"msg":"hello"}`
	got := decode(t, tr.Apply(line))
	if got["msg"] != "hello" {
		t.Errorf("unexpected change")
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	tr := New([]Rule{{Field: "count", MaxLen: 1}})
	line := `{"count":42}`
	got := decode(t, tr.Apply(line))
	if got["count"] != float64(42) {
		t.Errorf("expected 42, got %v", got["count"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 3}})
	line := `not json`
	if got := tr.Apply(line); got != line {
		t.Errorf("expected original line back")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	tr := New([]Rule{
		{Field: "msg", MaxLen: 3},
		{Field: "src", MaxLen: 4},
	})
	line := `{"msg":"hello","src":"server01"}`
	got := decode(t, tr.Apply(line))
	if got["msg"] != "hel" {
		t.Errorf("msg: expected 'hel', got %v", got["msg"])
	}
	if got["src"] != "serv" {
		t.Errorf("src: expected 'serv', got %v", got["src"])
	}
}
