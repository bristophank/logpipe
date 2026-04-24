package dedenter

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	d := New(nil)
	in := `{"msg":"  hello"}`
	if got := d.Apply(in); got != in {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_StripLeadingSpaces(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	out := d.Apply(`{"msg":"   hello world"}`)
	m := decode(t, out)
	if m["msg"] != "hello world" {
		t.Errorf("unexpected value: %v", m["msg"])
	}
}

func TestApply_StripLeadingTabs(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	out := d.Apply("{\"msg\":\"\\t\\thello\"}")
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("unexpected value: %v", m["msg"])
	}
}

func TestApply_CollapseInner(t *testing.T) {
	d := New([]Rule{{Field: "msg", CollapseInner: true}})
	out := d.Apply(`{"msg":"  hello   world  "}`)
	m := decode(t, out)
	if m["msg"] != "hello world" {
		t.Errorf("unexpected value: %v", m["msg"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	d := New([]Rule{{Field: "missing"}})
	in := `{"msg":"  hello"}`
	out := d.Apply(in)
	m := decode(t, out)
	if m["msg"] != "  hello" {
		t.Errorf("field should be unchanged: %v", m["msg"])
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	d := New([]Rule{{Field: "count"}})
	in := `{"count":42}`
	out := d.Apply(in)
	m := decode(t, out)
	if int(m["count"].(float64)) != 42 {
		t.Errorf("non-string field should be unchanged: %v", m["count"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	in := `not-json`
	if got := d.Apply(in); got != in {
		t.Errorf("expected passthrough for invalid JSON, got %q", got)
	}
}

func TestApply_MultilineField(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	in := "{\"msg\":\"  line1\\n  line2\"}" 
	out := d.Apply(in)
	m := decode(t, out)
	if m["msg"] != "line1\nline2" {
		t.Errorf("unexpected multiline result: %v", m["msg"])
	}
}
