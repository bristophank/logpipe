package caster

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return obj
}

func TestApply_NoRules(t *testing.T) {
	c := New(nil)
	input := `{"count":"42"}`
	if got := c.Apply(input); got != input {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestApply_CastToInt(t *testing.T) {
	c := New([]Rule{{Field: "count", Format: "int"}})
	out := c.Apply(`{"count":"7"}`)
	obj := decode(t, out)
	if obj["count"] != float64(7) {
		t.Errorf("expected 7, got %v", obj["count"])
	}
}

func TestApply_CastToFloat(t *testing.T) {
	c := New([]Rule{{Field: "ratio", Format: "float"}})
	out := c.Apply(`{"ratio":"3.14"}`)
	obj := decode(t, out)
	if obj["ratio"] != 3.14 {
		t.Errorf("expected 3.14, got %v", obj["ratio"])
	}
}

func TestApply_CastToBool(t *testing.T) {
	c := New([]Rule{{Field: "active", Format: "bool"}})
	out := c.Apply(`{"active":"true"}`)
	obj := decode(t, out)
	if obj["active"] != true {
		t.Errorf("expected true, got %v", obj["active"])
	}
}

func TestApply_CastToString(t *testing.T) {
	c := New([]Rule{{Field: "code", Format: "string"}})
	out := c.Apply(`{"code":404}`)
	obj := decode(t, out)
	if obj["code"] != "404" {
		t.Errorf("expected \"404\", got %v", obj["code"])
	}
}

func TestApply_MissingField_Unchanged(t *testing.T) {
	c := New([]Rule{{Field: "missing", Format: "int"}})
	input := `{"other":"val"}`
	out := c.Apply(input)
	obj := decode(t, out)
	if _, ok := obj["missing"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	c := New([]Rule{{Field: "x", Format: "int"}})
	input := `not-json`
	if got := c.Apply(input); got != input {
		t.Errorf("expected passthrough for invalid JSON, got %s", got)
	}
}

func TestApply_UnknownFormat_FieldUnchanged(t *testing.T) {
	c := New([]Rule{{Field: "val", Format: "hex"}})
	input := `{"val":"ff"}`
	out := c.Apply(input)
	obj := decode(t, out)
	if obj["val"] != "ff" {
		t.Errorf("expected original value preserved, got %v", obj["val"])
	}
}
