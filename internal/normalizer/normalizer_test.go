package normalizer

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
	n := New(nil)
	in := `{"level":"INFO","msg":"hello"}`
	if got := n.Apply(in); got != in {
		t.Fatalf("expected unchanged, got %s", got)
	}
}

func TestApply_Lowercase(t *testing.T) {
	n := New([]Rule{{Field: "level", Transform: "lowercase"}})
	out := decode(t, n.Apply(`{"level":"ERROR"}`))
	if out["level"] != "error" {
		t.Fatalf("expected 'error', got %v", out["level"])
	}
}

func TestApply_Uppercase(t *testing.T) {
	n := New([]Rule{{Field: "env", Transform: "uppercase"}})
	out := decode(t, n.Apply(`{"env":"production"}`))
	if out["env"] != "PRODUCTION" {
		t.Fatalf("expected 'PRODUCTION', got %v", out["env"])
	}
}

func TestApply_Trim(t *testing.T) {
	n := New([]Rule{{Field: "msg", Transform: "trim"}})
	out := decode(t, n.Apply(`{"msg":"  hello  "}`))
	if out["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %v", out["msg"])
	}
}

func TestApply_SnakeCase(t *testing.T) {
	n := New([]Rule{{Field: "service", Transform: "snake_case"}})
	out := decode(t, n.Apply(`{"service":"My-Service Name"}`))
	if out["service"] != "my_service_name" {
		t.Fatalf("expected 'my_service_name', got %v", out["service"])
	}
}

func TestApply_Rename(t *testing.T) {
	n := New([]Rule{{Field: "level", Transform: "lowercase", Rename: "severity"}})
	out := decode(t, n.Apply(`{"level":"WARN"}`))
	if _, ok := out["level"]; ok {
		t.Fatal("old field 'level' should be removed")
	}
	if out["severity"] != "warn" {
		t.Fatalf("expected 'warn', got %v", out["severity"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	n := New([]Rule{{Field: "missing", Transform: "lowercase"}})
	in := `{"level":"INFO"}`
	out := n.Apply(in)
	if decode(t, out)["level"] != "INFO" {
		t.Fatal("unrelated field should be unchanged")
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	n := New([]Rule{{Field: "level", Transform: "lowercase"}})
	in := `not json`
	if got := n.Apply(in); got != in {
		t.Fatalf("expected passthrough, got %s", got)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	n := New([]Rule{{Field: "count", Transform: "lowercase"}})
	in := `{"count":42}`
	out := decode(t, n.Apply(in))
	if out["count"] != float64(42) {
		t.Fatalf("numeric field should be unchanged, got %v", out["count"])
	}
}
