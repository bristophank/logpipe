package masker

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
	m := New(nil)
	in := `{"password":"secret"}`
	if got := m.Apply(in); got != in {
		t.Fatalf("expected unchanged, got %s", got)
	}
}

func TestApply_DefaultMask(t *testing.T) {
	m := New([]Rule{{Field: "password"}})
	out := m.Apply(`{"password":"hunter2"}`)
	obj := decode(t, out)
	if obj["password"] != "***" {
		t.Fatalf("expected ***, got %v", obj["password"])
	}
}

func TestApply_CustomMask(t *testing.T) {
	m := New([]Rule{{Field: "token", Mask: "REDACTED"}})
	out := m.Apply(`{"token":"abc123"}`)
	obj := decode(t, out)
	if obj["token"] != "REDACTED" {
		t.Fatalf("expected REDACTED, got %v", obj["token"])
	}
}

func TestApply_KeepTrailing(t *testing.T) {
	m := New([]Rule{{Field: "card", Keep: 4}})
	out := m.Apply(`{"card":"1234567890123456"}`)
	obj := decode(t, out)
	if obj["card"] != "***3456" {
		t.Fatalf("unexpected value: %v", obj["card"])
	}
}

func TestApply_MissingField(t *testing.T) {
	m := New([]Rule{{Field: "secret"}})
	in := `{"level":"info"}`
	out := m.Apply(in)
	obj := decode(t, out)
	if _, ok := obj["secret"]; ok {
		t.Fatal("field should not exist")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m := New([]Rule{{Field: "x"}})
	in := "not-json"
	if got := m.Apply(in); got != in {
		t.Fatalf("expected passthrough, got %s", got)
	}
}

func TestApply_NonStringField(t *testing.T) {
	m := New([]Rule{{Field: "count"}})
	in := `{"count":42}`
	out := m.Apply(in)
	obj := decode(t, out)
	if obj["count"] != float64(42) {
		t.Fatalf("non-string field should be unchanged: %v", obj["count"])
	}
}
