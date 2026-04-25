package bouncer

import (
	"testing"
)

func TestAllow_NoRules(t *testing.T) {
	b := New(nil)
	if !b.Allow(`{"level":"error"}`) {
		t.Fatal("expected allow with no rules")
	}
}

func TestAllow_InvalidJSON_Passes(t *testing.T) {
	b := New([]Rule{{Field: "level", Allow: []string{"info"}}})
	if !b.Allow("not-json") {
		t.Fatal("invalid JSON should pass through")
	}
}

func TestAllow_AllowlistMatch(t *testing.T) {
	b := New([]Rule{{Field: "level", Allow: []string{"info", "warn"}}})
	if !b.Allow(`{"level":"info","msg":"ok"}`) {
		t.Fatal("expected allow for matching allowlist value")
	}
}

func TestAllow_AllowlistNoMatch(t *testing.T) {
	b := New([]Rule{{Field: "level", Allow: []string{"info", "warn"}}})
	if b.Allow(`{"level":"error","msg":"boom"}`) {
		t.Fatal("expected drop for non-allowlist value")
	}
}

func TestAllow_BlocklistMatch(t *testing.T) {
	b := New([]Rule{{Field: "level", Block: []string{"debug"}}})
	if b.Allow(`{"level":"debug","msg":"verbose"}`) {
		t.Fatal("expected drop for blocked value")
	}
}

func TestAllow_BlocklistNoMatch(t *testing.T) {
	b := New([]Rule{{Field: "level", Block: []string{"debug"}}})
	if !b.Allow(`{"level":"info","msg":"ok"}`) {
		t.Fatal("expected allow for non-blocked value")
	}
}

func TestAllow_MissingField_Passes(t *testing.T) {
	b := New([]Rule{{Field: "env", Allow: []string{"prod"}}})
	if !b.Allow(`{"level":"info"}`) {
		t.Fatal("missing field should not be filtered")
	}
}

func TestAllow_MultipleRules_AllMustPass(t *testing.T) {
	b := New([]Rule{
		{Field: "level", Allow: []string{"info"}},
		{Field: "env", Block: []string{"staging"}},
	})
	if b.Allow(`{"level":"info","env":"staging"}`) {
		t.Fatal("expected drop: env is blocked")
	}
	if !b.Allow(`{"level":"info","env":"prod"}`) {
		t.Fatal("expected allow: both rules pass")
	}
}
