package alerter

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_PassesAllLines(t *testing.T) {
	a := New(nil)
	input := "{\"msg\":\"hello\"}\n{\"msg\":\"world\"}\n"
	var out, alerts bytes.Buffer
	if err := a.Stream(strings.NewReader(input), &out, &alerts); err != nil {
		t.Fatalf("stream error: %v", err)
	}
	if strings.Count(out.String(), "\n") != 2 {
		t.Fatalf("expected 2 lines, got: %q", out.String())
	}
	if alerts.Len() != 0 {
		t.Fatalf("expected no alerts")
	}
}

func TestStream_WritesAlerts(t *testing.T) {
	a := New([]Rule{{Field: "status", Threshold: 400, Window: 60}})
	input := "{\"status\":500}\n{\"status\":200}\n"
	var out, alerts bytes.Buffer
	if err := a.Stream(strings.NewReader(input), &out, &alerts); err != nil {
		t.Fatalf("stream error: %v", err)
	}
	if !strings.Contains(alerts.String(), `"alert":true`) {
		t.Fatalf("expected alert, got: %q", alerts.String())
	}
	if strings.Count(out.String(), "\n") != 2 {
		t.Fatalf("all lines should pass through")
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	a := New(nil)
	input := "{\"a\":1}\n\n{\"b\":2}\n"
	var out, alerts bytes.Buffer
	_ = a.Stream(strings.NewReader(input), &out, &alerts)
	if strings.Count(out.String(), "\n") != 2 {
		t.Fatalf("expected 2 non-empty lines")
	}
}
