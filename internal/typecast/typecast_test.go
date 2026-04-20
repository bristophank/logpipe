package typecast

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
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c := New([]Rule{{Field: "x", To: "int"}})
	out, err := c.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
	if out != `not json` {
		t.Fatalf("expected original line returned, got: %s", out)
	}
}

func TestApply_CastToInt(t *testing.T) {
	c := New([]Rule{{Field: "code", To: "int"}})
	out, err := c.Apply(`{"code":"42"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if v, ok := m["code"].(float64); !ok || v != 42 {
		t.Fatalf("expected code=42 (float64 from JSON), got %v (%T)", m["code"], m["code"])
	}
}

func TestApply_CastToFloat(t *testing.T) {
	c := New([]Rule{{Field: "latency", To: "float"}})
	out, err := c.Apply(`{"latency":"1.23"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if v, ok := m["latency"].(float64); !ok || v != 1.23 {
		t.Fatalf("expected latency=1.23, got %v", m["latency"])
	}
}

func TestApply_CastToBool(t *testing.T) {
	c := New([]Rule{{Field: "ok", To: "bool"}})
	out, err := c.Apply(`{"ok":"true"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if v, ok := m["ok"].(bool); !ok || !v {
		t.Fatalf("expected ok=true, got %v", m["ok"])
	}
}

func TestApply_CastToString(t *testing.T) {
	c := New([]Rule{{Field: "status", To: "string"}})
	out, err := c.Apply(`{"status":200}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if v, ok := m["status"].(string); !ok || v != "200" {
		t.Fatalf("expected status=\"200\", got %v (%T)", m["status"], m["status"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	c := New([]Rule{{Field: "missing", To: "int"}})
	out, err := c.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, exists := m["missing"]; exists {
		t.Fatal("missing field should not be added")
	}
}

func TestApply_InvalidCast_Skipped(t *testing.T) {
	c := New([]Rule{{Field: "val", To: "int"}})
	out, err := c.Apply(`{"val":"notanumber"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	// original string value should remain unchanged
	if v, ok := m["val"].(string); !ok || v != "notanumber" {
		t.Fatalf("expected original value preserved, got %v", m["val"])
	}
}
