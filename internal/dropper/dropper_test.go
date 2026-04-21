package dropper

import (
	"testing"
)

func TestShouldDrop_NoRules(t *testing.T) {
	d := New(nil)
	if d.ShouldDrop(`{"level":"error"}`) {
		t.Fatal("expected no drop with no rules")
	}
}

func TestShouldDrop_EqMatch(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	if !d.ShouldDrop(`{"level":"debug","msg":"verbose"}`) {
		t.Fatal("expected line to be dropped")
	}
}

func TestShouldDrop_EqNoMatch(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	if d.ShouldDrop(`{"level":"error","msg":"oops"}`) {
		t.Fatal("expected line to pass through")
	}
}

func TestShouldDrop_ContainsMatch(t *testing.T) {
	d := New([]Rule{{Field: "msg", Operator: "contains", Value: "healthcheck"}})
	if !d.ShouldDrop(`{"msg":"GET /healthcheck 200"}`) {
		t.Fatal("expected healthcheck line to be dropped")
	}
}

func TestShouldDrop_ExistsMatch(t *testing.T) {
	d := New([]Rule{{Field: "internal", Operator: "exists"}})
	if !d.ShouldDrop(`{"internal":true,"msg":"skip me"}`) {
		t.Fatal("expected line with 'internal' field to be dropped")
	}
}

func TestShouldDrop_ExistsNoMatch(t *testing.T) {
	d := New([]Rule{{Field: "internal", Operator: "exists"}})
	if d.ShouldDrop(`{"msg":"keep me"}`) {
		t.Fatal("expected line without 'internal' field to pass")
	}
}

func TestShouldDrop_InvalidJSON(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	if d.ShouldDrop(`not json`) {
		t.Fatal("expected invalid JSON to pass through")
	}
}

func TestShouldDrop_EmptyLine(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	if d.ShouldDrop("") {
		t.Fatal("expected empty line to pass through")
	}
}

func TestApply_DropsLine(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	result := d.Apply(`{"level":"debug"}`)
	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

func TestApply_PassesLine(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	line := `{"level":"info","msg":"hello"}`
	result := d.Apply(line)
	if result != line {
		t.Fatalf("expected line unchanged, got %q", result)
	}
}
