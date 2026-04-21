package limiter

import (
	"testing"
)

func TestAllow_NoRules(t *testing.T) {
	l := New(nil)
	if !l.Allow(`{"msg":"hello world this is a very long message"}`) {
		t.Fatal("expected allow with no rules")
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 20}})
	if !l.Allow(`{"msg":"short"}`) {
		t.Fatal("expected allow: value within limit")
	}
}

func TestAllow_ExactLimit(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	if !l.Allow(`{"msg":"hello"}`) {
		t.Fatal("expected allow: value exactly at limit")
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	if l.Allow(`{"msg":"toolong"}`) {
		t.Fatal("expected drop: value exceeds limit")
	}
}

func TestAllow_MissingField_Passes(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	if !l.Allow(`{"level":"info"}`) {
		t.Fatal("expected allow: field absent")
	}
}

func TestAllow_NonStringField_Passes(t *testing.T) {
	l := New([]Rule{{Field: "count", MaxLen: 2}})
	if !l.Allow(`{"count":12345}`) {
		t.Fatal("expected allow: non-string field not evaluated")
	}
}

func TestAllow_InvalidJSON_Passes(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	if !l.Allow(`not json`) {
		t.Fatal("expected allow: invalid JSON passes through")
	}
}

func TestAllow_EmptyLine_Passes(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	if !l.Allow("") {
		t.Fatal("expected allow: empty line passes through")
	}
}

func TestAllow_MultipleRules_AnyExceeds(t *testing.T) {
	l := New([]Rule{
		{Field: "msg", MaxLen: 20},
		{Field: "trace", MaxLen: 5},
	})
	line := `{"msg":"ok","trace":"toolongvalue"}`
	if l.Allow(line) {
		t.Fatal("expected drop: second rule exceeded")
	}
}

func TestAllow_ZeroMaxLen_RuleIgnored(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 0}})
	if !l.Allow(`{"msg":"anything goes here"}`) {
		t.Fatal("expected allow: zero MaxLen rule is ignored")
	}
}
