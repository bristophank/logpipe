package timestamper

import (
	"encoding/json"
	"testing"
	"time"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func fixed() func() time.Time { return func() time.Time { return fixedTime } }

func TestApply_NoRules(t *testing.T) {
	ts := New(nil, fixed())
	out, err := ts.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"msg":"hello"}` {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestApply_AddsTimestamp(t *testing.T) {
	ts := New([]Rule{{Field: "ts"}}, fixed())
	out, err := ts.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["ts"] != fixedTime.UTC().Format(time.RFC3339) {
		t.Errorf("unexpected ts: %v", m["ts"])
	}
}

func TestApply_NoOverwrite(t *testing.T) {
	ts := New([]Rule{{Field: "ts", Overwrite: false}}, fixed())
	out, err := ts.Apply(`{"ts":"existing"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["ts"] != "existing" {
		t.Errorf("expected field preserved, got: %v", m["ts"])
	}
}

func TestApply_Overwrite(t *testing.T) {
	ts := New([]Rule{{Field: "ts", Overwrite: true}}, fixed())
	out, err := ts.Apply(`{"ts":"old"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["ts"] == "old" {
		t.Error("expected field to be overwritten")
	}
}

func TestApply_CustomFormat(t *testing.T) {
	format := "2006-01-02"
	ts := New([]Rule{{Field: "date", Format: format}}, fixed())
	out, err := ts.Apply(`{"msg":"x"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["date"] != fixedTime.UTC().Format(format) {
		t.Errorf("unexpected date: %v", m["date"])
	}
}

func TestApply_InvalidJSON_Passthrough(t *testing.T) {
	ts := New([]Rule{{Field: "ts"}}, fixed())
	raw := `not json`
	out, err := ts.Apply(raw)
	if err != nil {
		t.Fatal(err)
	}
	if out != raw {
		t.Errorf("expected passthrough, got: %s", out)
	}
}
