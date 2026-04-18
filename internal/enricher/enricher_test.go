package enricher

import (
	"encoding/json"
	"strings"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	e := New(nil, "host1")
	in := `{"level":"info"}`
	if got := e.Apply(in); got != in {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestApply_StaticValue(t *testing.T) {
	e := New([]Rule{{Field: "env", Value: "prod"}}, "")
	out := decode(t, e.Apply(`{"msg":"ok"}`))
	if out["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", out["env"])
	}
}

func TestApply_HostnameSource(t *testing.T) {
	e := New([]Rule{{Field: "host", Source: "hostname"}}, "myhost")
	out := decode(t, e.Apply(`{"msg":"ok"}`))
	if out["host"] != "myhost" {
		t.Errorf("expected host=myhost, got %v", out["host"])
	}
}

func TestApply_TimestampSource(t *testing.T) {
	e := New([]Rule{{Field: "ts", Source: "timestamp"}}, "")
	out := decode(t, e.Apply(`{"msg":"ok"}`))
	ts, ok := out["ts"].(string)
	if !ok || ts == "" {
		t.Errorf("expected non-empty ts, got %v", out["ts"])
	}
}

func TestApply_NonJSON_Passthrough(t *testing.T) {
	e := New([]Rule{{Field: "env", Value: "prod"}}, "")
	in := "plain text line"
	if got := e.Apply(in); got != in {
		t.Errorf("expected passthrough for non-json, got %s", got)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	rules := []Rule{
		{Field: "env", Value: "staging"},
		{Field: "host", Source: "hostname"},
	}
	e := New(rules, "srv1")
	out := decode(t, e.Apply(`{"level":"warn"}`))
	if out["env"] != "staging" || out["host"] != "srv1" {
		t.Errorf("unexpected fields: %v", out)
	}
	if !strings.Contains(e.Apply(`{"x":1}`), "env") {
		t.Error("expected env field in output")
	}
}
