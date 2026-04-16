package filter

import "testing"

func TestMatch_NoRules(t *testing.T) {
	f := New(nil)
	if !f.Match(`{"level":"info","msg":"hello"}`) {
		t.Error("expected match with no rules")
	}
}

func TestMatch_EqRule(t *testing.T) {
	f := New([]Rule{{Field: "level", Operator: "eq", Value: "error"}})

	if f.Match(`{"level":"info","msg":"ok"}`) {
		t.Error("expected no match for level=info")
	}
	if !f.Match(`{"level":"error","msg":"fail"}`) {
		t.Error("expected match for level=error")
	}
}

func TestMatch_ContainsRule(t *testing.T) {
	f := New([]Rule{{Field: "msg", Operator: "contains", Value: "timeout"}})

	if !f.Match(`{"msg":"connection timeout reached"}`) {
		t.Error("expected match containing 'timeout'")
	}
	if f.Match(`{"msg":"all good"}`) {
		t.Error("expected no match")
	}
}

func TestMatch_ExistsRule(t *testing.T) {
	f := New([]Rule{{Field: "trace_id", Operator: "exists"}})

	if !f.Match(`{"trace_id":"abc123","msg":"ok"}`) {
		t.Error("expected match when field exists")
	}
	if f.Match(`{"msg":"no trace"}`) {
		t.Error("expected no match when field absent")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f := New([]Rule{{Field: "level", Operator: "eq", Value: "info"}})
	if f.Match(`not json`) {
		t.Error("expected no match for invalid JSON")
	}
}

func TestMatch_MultipleRules(t *testing.T) {
	f := New([]Rule{
		{Field: "level", Operator: "eq", Value: "error"},
		{Field: "service", Operator: "contains", Value: "auth"},
	})

	if !f.Match(`{"level":"error","service":"auth-service"}`) {
		t.Error("expected match for both rules satisfied")
	}
	if f.Match(`{"level":"error","service":"billing"}`) {
		t.Error("expected no match when second rule fails")
	}
}
