package shifter

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
	sh := New(nil)
	line := `{"count":5}`
	got, err := sh.Apply(line)
	if err != nil {
		t.Fatal(err)
	}
	if got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	sh := New([]Rule{{Field: "x", By: 1}})
	got, err := sh.Apply("not-json")
	if err != nil {
		t.Fatal(err)
	}
	if got != "not-json" {
		t.Errorf("expected passthrough for invalid JSON, got %q", got)
	}
}

func TestApply_ShiftByPositive(t *testing.T) {
	sh := New([]Rule{{Field: "count", By: 10}})
	got, err := sh.Apply(`{"count":5}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, got)
	if m["count"].(float64) != 15 {
		t.Errorf("expected 15, got %v", m["count"])
	}
}

func TestApply_ShiftByNegative(t *testing.T) {
	sh := New([]Rule{{Field: "score", By: -3}})
	got, _ := sh.Apply(`{"score":10}`)
	m := decode(t, got)
	if m["score"].(float64) != 7 {
		t.Errorf("expected 7, got %v", m["score"])
	}
}

func TestApply_WithScale(t *testing.T) {
	sh := New([]Rule{{Field: "val", By: 0, Scale: 2.5}})
	got, _ := sh.Apply(`{"val":4}`)
	m := decode(t, got)
	if m["val"].(float64) != 10 {
		t.Errorf("expected 10, got %v", m["val"])
	}
}

func TestApply_TargetField(t *testing.T) {
	sh := New([]Rule{{Field: "price", By: 5, Target: "adjusted"}})
	got, _ := sh.Apply(`{"price":20}`)
	m := decode(t, got)
	if m["adjusted"].(float64) != 25 {
		t.Errorf("expected adjusted=25, got %v", m["adjusted"])
	}
	if m["price"].(float64) != 20 {
		t.Errorf("expected original price unchanged, got %v", m["price"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	sh := New([]Rule{{Field: "missing", By: 1}})
	line := `{"other":9}`
	got, _ := sh.Apply(line)
	m := decode(t, got)
	if _, ok := m["missing"]; ok {
		t.Error("missing field should not be created")
	}
}

func TestApply_DefaultScaleIsOne(t *testing.T) {
	sh := New([]Rule{{Field: "n", By: 0, Scale: 0}})
	got, _ := sh.Apply(`{"n":7}`)
	m := decode(t, got)
	if m["n"].(float64) != 7 {
		t.Errorf("expected 7 (scale=1), got %v", m["n"])
	}
}
