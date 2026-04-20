package joiner

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
	j := New(nil)
	line := `{"user_id":"1","action":"login"}`
	out, err := j.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Fatalf("expected passthrough, got %s", out)
	}
}

func TestApply_JoinsFields(t *testing.T) {
	rules := []Rule{{PrimaryKey: "user_id", SecondaryKey: "id", Fields: []string{"name", "email"}}}
	j := New(rules)

	_ = j.Index(0, `{"id":"42","name":"alice","email":"alice@example.com"}`)

	out, err := j.Apply(`{"user_id":"42","action":"login"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["name"] != "alice" {
		t.Errorf("expected name=alice, got %v", m["name"])
	}
	if m["email"] != "alice@example.com" {
		t.Errorf("expected email set, got %v", m["email"])
	}
	if m["action"] != "login" {
		t.Errorf("expected action preserved, got %v", m["action"])
	}
}

func TestApply_NoMatchingKey(t *testing.T) {
	rules := []Rule{{PrimaryKey: "user_id", SecondaryKey: "id", Fields: []string{"name"}}}
	j := New(rules)
	_ = j.Index(0, `{"id":"99","name":"bob"}`)

	out, err := j.Apply(`{"user_id":"1","action":"logout"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["name"]; ok {
		t.Error("expected name not to be joined")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	j := New([]Rule{{PrimaryKey: "id", SecondaryKey: "id", Fields: []string{"x"}}})
	out, err := j.Apply("not-json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "not-json" {
		t.Errorf("expected passthrough for invalid JSON, got %s", out)
	}
}

func TestIndex_InvalidRuleIndex(t *testing.T) {
	j := New(nil)
	err := j.Index(5, `{"id":"1"}`)
	if err == nil {
		t.Fatal("expected error for out-of-range rule index")
	}
}

func TestReset_ClearsTable(t *testing.T) {
	rules := []Rule{{PrimaryKey: "user_id", SecondaryKey: "id", Fields: []string{"name"}}}
	j := New(rules)
	_ = j.Index(0, `{"id":"1","name":"alice"}`)
	j.Reset()

	out, _ := j.Apply(`{"user_id":"1"}`)
	m := decode(t, out)
	if _, ok := m["name"]; ok {
		t.Error("expected name not present after reset")
	}
}
