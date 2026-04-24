package unwrapper

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_UnwrapsLines(t *testing.T) {
	u := New([]Rule{{Field: "meta", Prefix: "m_", Delete: true}})
	input := `{"level":"info","meta":{"host":"a"}}
{"level":"warn","meta":{"host":"b"}}
`
	var buf bytes.Buffer
	if err := Stream(u, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		m := decode(t, l)
		if _, ok := m["m_host"]; !ok {
			t.Errorf("expected m_host in output line: %s", l)
		}
		if _, ok := m["meta"]; ok {
			t.Errorf("expected meta deleted in output line: %s", l)
		}
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	u := New([]Rule{{Field: "meta", Delete: true}})
	input := `{"level":"info","meta":{"host":"x"}}

{"level":"debug","meta":{"host":"y"}}
`
	var buf bytes.Buffer
	if err := Stream(u, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	u := New([]Rule{{Field: "meta", Delete: true}})
	input := "not-json\n"
	var buf bytes.Buffer
	if err := Stream(u, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	// invalid JSON is passed through as-is by Apply, so it should appear
	if !strings.Contains(buf.String(), "not-json") {
		t.Errorf("expected invalid JSON to pass through, got: %s", buf.String())
	}
}
