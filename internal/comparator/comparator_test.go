package comparator

import (
	"encoding/json"
	"testing"
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
	c := New(nil)
	out, err := c.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"level":"info"}` {
		t.Fatalf("expected passthrough, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c := New([]Rule{{Field: "x", Op: "gt", Value: 1, Target: "result"}})
	out, err := c.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error")
	}
	if out != `not-json` {
		t.Fatalf("expected original line on error, got %s", out)
	}
}

func TestApply_GtTrue(t *testing.T) {
	c := New([]Rule{{Field: "count", Op: "gt", Value: 5, Target: "is_high"}})
	out, err := c.Apply(`{"count":10}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["is_high"] != true {
		t.Fatalf("expected true, got %v", m["is_high"])
	}
}

func TestApply_GtFalse(t *testing.T) {
	c := New([]Rule{{Field: "count", Op: "gt", Value: 5, Target: "is_high"}})
	out, err := c.Apply(`{"count":3}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["is_high"] != false {
		t.Fatalf("expected false, got %v", m["is_high"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	c := New([]Rule{{Field: "missing", Op: "lt", Value: 10, Target: "result"}})
	out, err := c.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["result"]; ok {
		t.Fatal("expected no result field for missing source")
	}
}

func TestApply_DefaultTargetName(t *testing.T) {
	c := New([]Rule{{Field: "score", Op: "gte", Value: 50}})
	out, err := c.Apply(`{"score":75}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["score_gte"] != true {
		t.Fatalf("expected score_gte=true, got %v", m["score_gte"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	c := New([]Rule{
		{Field: "latency", Op: "gt", Value: 100, Target: "slow"},
		{Field: "latency", Op: "lte", Value: 500, Target: "acceptable"},
	})
	out, err := c.Apply(`{"latency":200}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["slow"] != true {
		t.Fatalf("expected slow=true")
	}
	if m["acceptable"] != true {
		t.Fatalf("expected acceptable=true")
	}
}
