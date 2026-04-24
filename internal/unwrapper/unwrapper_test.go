package unwrapper

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
	u := New(nil)
	in := `{"level":"info","meta":{"host":"srv1"}}`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_UnwrapsNestedObject(t *testing.T) {
	u := New([]Rule{{Field: "meta", Delete: true}})
	in := `{"level":"info","meta":{"host":"srv1","env":"prod"}}`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %v", m["host"])
	}
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", m["env"])
	}
	if _, exists := m["meta"]; exists {
		t.Error("expected meta to be deleted")
	}
}

func TestApply_WithPrefix(t *testing.T) {
	u := New([]Rule{{Field: "meta", Prefix: "meta_", Delete: true}})
	in := `{"level":"warn","meta":{"host":"srv2"}}`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["meta_host"] != "srv2" {
		t.Errorf("expected meta_host=srv2, got %v", m["meta_host"])
	}
}

func TestApply_KeepsOriginalWhenDeleteFalse(t *testing.T) {
	u := New([]Rule{{Field: "meta", Delete: false}})
	in := `{"level":"debug","meta":{"host":"srv3"}}`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, exists := m["meta"]; !exists {
		t.Error("expected meta to be retained")
	}
	if m["host"] != "srv3" {
		t.Errorf("expected host=srv3, got %v", m["host"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	u := New([]Rule{{Field: "ctx", Delete: true}})
	in := `{"level":"info","msg":"hello"}`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", m["msg"])
	}
}

func TestApply_NonObjectField_Unchanged(t *testing.T) {
	u := New([]Rule{{Field: "level", Delete: true}})
	in := `{"level":"info","msg":"test"}`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	u := New([]Rule{{Field: "meta", Delete: true}})
	in := `not-json`
	out, err := u.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Errorf("expected passthrough for invalid JSON")
	}
}
