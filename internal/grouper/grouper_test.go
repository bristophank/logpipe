package grouper_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/grouper"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func collect(t *testing.T, g *grouper.Grouper) []map[string]any {
	t.Helper()
	var results []map[string]any
	for _, line := range g.Flush() {
		results = append(results, decode(t, line))
	}
	return results
}

func TestAdd_NoRules_PassesThrough(t *testing.T) {
	g := grouper.New(nil, 0)
	lines := []string{
		`{"level":"info","msg":"hello"}`,
		`{"level":"error","msg":"boom"}`,
	}
	for _, l := range lines {
		g.Add(l)
	}
	results := collect(t, g)
	if len(results) != 0 {
		t.Fatalf("expected 0 groups with no rules, got %d", len(results))
	}
}

func TestAdd_InvalidJSON_Skipped(t *testing.T) {
	rules := []grouper.Rule{{Field: "level"}}
	g := grouper.New(rules, 0)
	g.Add("not-json")
	results := collect(t, g)
	if len(results) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(results))
	}
}

func TestFlush_GroupsByField(t *testing.T) {
	rules := []grouper.Rule{{Field: "level"}}
	g := grouper.New(rules, 0)

	g.Add(`{"level":"info","msg":"a"}`)
	g.Add(`{"level":"error","msg":"b"}`)
	g.Add(`{"level":"info","msg":"c"}`)

	results := collect(t, g)
	if len(results) != 2 {
		t.Fatalf("expected 2 groups (info, error), got %d", len(results))
	}

	counts := map[string]int{}
	for _, r := range results {
		key, _ := r["level"].(string)
		cnt, _ := r["count"].(float64)
		counts[key] = int(cnt)
	}
	if counts["info"] != 2 {
		t.Errorf("expected info count=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error count=1, got %d", counts["error"])
	}
}

func TestFlush_ResetsState(t *testing.T) {
	rules := []grouper.Rule{{Field: "level"}}
	g := grouper.New(rules, 0)

	g.Add(`{"level":"info","msg":"a"}`)
	first := collect(t, g)
	if len(first) != 1 {
		t.Fatalf("expected 1 group after first flush, got %d", len(first))
	}

	second := collect(t, g)
	if len(second) != 0 {
		t.Fatalf("expected 0 groups after second flush (reset), got %d", len(second))
	}
}

func TestFlush_MissingField_GroupedUnderEmpty(t *testing.T) {
	rules := []grouper.Rule{{Field: "level"}}
	g := grouper.New(rules, 0)

	g.Add(`{"msg":"no level here"}`)
	g.Add(`{"msg":"also no level"}`)

	results := collect(t, g)
	if len(results) != 1 {
		t.Fatalf("expected 1 group for missing field, got %d", len(results))
	}
	cnt, _ := results[0]["count"].(float64)
	if int(cnt) != 2 {
		t.Errorf("expected count=2, got %d", int(cnt))
	}
}

func TestFlush_WindowExpiry(t *testing.T) {
	rules := []grouper.Rule{{Field: "level"}}
	window := 50 * time.Millisecond
	g := grouper.New(rules, window)

	g.Add(`{"level":"info","msg":"a"}`)
	time.Sleep(window + 20*time.Millisecond)

	// After window expiry, flush should return the accumulated group
	results := collect(t, g)
	if len(results) != 1 {
		t.Fatalf("expected 1 group after window expiry, got %d", len(results))
	}
}

func TestStream_GroupsAndFlushes(t *testing.T) {
	rules := []grouper.Rule{{Field: "svc"}}
	g := grouper.New(rules, 0)

	input := strings.NewReader(
		`{"svc":"api","msg":"req"}` + "\n" +
			`{"svc":"db","msg":"query"}` + "\n" +
			`{"svc":"api","msg":"resp"}` + "\n",
	)

	var out strings.Builder
	grouper.Stream(g, input, &out)

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d: %v", len(lines), lines)
	}
}
