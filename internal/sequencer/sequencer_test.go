package sequencer

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	seq := New(nil)
	out, err := seq.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `{"msg":"hello"}` {
		t.Fatalf("expected passthrough, got %q", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	seq := New([]Rule{{Field: "seq"}})
	out, err := seq.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if out != `not-json` {
		t.Fatalf("expected original line on error, got %q", out)
	}
}

func TestApply_DefaultStartAndStep(t *testing.T) {
	seq := New([]Rule{{Field: "n"}})
	for i := 0; i < 3; i++ {
		out, err := seq.Apply(`{"x":1}`)
		if err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
		m := decode(t, out)
		got := int(m["n"].(float64))
		if got != i {
			t.Fatalf("step %d: expected n=%d, got %d", i, i, got)
		}
	}
}

func TestApply_CustomStartAndStep(t *testing.T) {
	seq := New([]Rule{{Field: "seq", Start: 10, Step: 5}})
	expected := []int{10, 15, 20}
	for i, exp := range expected {
		out, _ := seq.Apply(`{"msg":"x"}`)
		m := decode(t, out)
		got := int(m["seq"].(float64))
		if got != exp {
			t.Fatalf("step %d: expected %d, got %d", i, exp, got)
		}
	}
}

func TestApply_MultipleRules(t *testing.T) {
	seq := New([]Rule{
		{Field: "a", Start: 0, Step: 1},
		{Field: "b", Start: 100, Step: 10},
	})
	out, _ := seq.Apply(`{"msg":"y"}`)
	m := decode(t, out)
	if int(m["a"].(float64)) != 0 {
		t.Fatalf("expected a=0")
	}
	if int(m["b"].(float64)) != 100 {
		t.Fatalf("expected b=100")
	}
}

func TestReset(t *testing.T) {
	seq := New([]Rule{{Field: "seq", Start: 1, Step: 1}})
	seq.Apply(`{"x":1}`)
	seq.Apply(`{"x":2}`)
	seq.Reset()
	out, _ := seq.Apply(`{"x":3}`)
	m := decode(t, out)
	if int(m["seq"].(float64)) != 1 {
		t.Fatalf("expected seq=1 after reset")
	}
}

func TestSnapshot(t *testing.T) {
	seq := New([]Rule{{Field: "idx", Start: 0, Step: 2}})
	seq.Apply(`{"a":1}`)
	seq.Apply(`{"a":2}`)
	snap := seq.Snapshot()
	if snap["idx"] != 4 {
		t.Fatalf("expected idx=4, got %d", snap["idx"])
	}
}
