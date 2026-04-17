package redactor

import (
	"encoding/json"
	"testing"
)

func TestApply_NoRules(t *testing.T) {
	r, _ := New(nil)
	line := `{"password":"secret"}`
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged, got %s", got)
	}
}

func TestApply_MaskWholeField(t *testing.T) {
	r, _ := New([]Rule{{Field: "password"}})
	got := r.Apply(`{"user":"alice","password":"s3cr3t"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(got), &obj); err != nil {
		t.Fatal(err)
	}
	if obj["password"] != "***" {
		t.Fatalf("expected ***, got %v", obj["password"])
	}
	if obj["user"] != "alice" {
		t.Fatalf("user field altered")
	}
}

func TestApply_CustomMask(t *testing.T) {
	r, _ := New([]Rule{{Field: "token", Mask: "[REDACTED]"}})
	got := r.Apply(`{"token":"abc123"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(got), &obj)
	if obj["token"] != "[REDACTED]" {
		t.Fatalf("unexpected mask: %v", obj["token"])
	}
}

func TestApply_PatternPartialRedact(t *testing.T) {
	r, _ := New([]Rule{{Field: "email", Pattern: `@.*`, Mask: "@***"}})
	got := r.Apply(`{"email":"user@example.com"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(got), &obj)
	if obj["email"] != "user@***" {
		t.Fatalf("unexpected value: %v", obj["email"])
	}
}

func TestApply_MissingField(t *testing.T) {
	r, _ := New([]Rule{{Field: "secret"}})
	line := `{"msg":"hello"}`
	if got := r.Apply(line); got != line {
		// field absent — object unchanged structurally, just re-encoded
		var a, b map[string]interface{}
		json.Unmarshal([]byte(line), &a)
		json.Unmarshal([]byte(got), &b)
		if len(a) != len(b) {
			t.Fatalf("unexpected change: %s", got)
		}
	}
}

func TestApply_NonJSON(t *testing.T) {
	r, _ := New([]Rule{{Field: "x"}})
	line := "not json"
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged, got %s", got)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]Rule{{Field: "f", Pattern: "[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}
