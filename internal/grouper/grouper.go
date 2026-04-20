// Package grouper collects log lines and groups them by a specified field value,
// emitting a summary JSON object for each group when flushed.
package grouper

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// Rule defines how lines should be grouped.
type Rule struct {
	// Field is the JSON key whose value determines the group.
	Field string `json:"field"`
	// OutputField is the key used in the emitted summary object (defaults to Field).
	OutputField string `json:"output_field"`
	// CountField is the key used to store the count in the summary (defaults to "count").
	CountField string `json:"count_field"`
}

// group holds accumulated lines for a single field value.
type group struct {
	value string
	lines []map[string]any
}

// Grouper batches incoming JSON log lines and groups them by a field value.
type Grouper struct {
	rules  []Rule
	mu     sync.Mutex
	// groups maps rule-index -> field-value -> accumulated entries
	groups map[int]map[string]*group
}

// New creates a Grouper with the provided rules.
func New(rules []Rule) *Grouper {
	g := &Grouper{
		rules:  rules,
		groups: make(map[int]map[string]*group),
	}
	for i := range rules {
		g.groups[i] = make(map[string]*group)
	}
	return g
}

// Add parses a JSON log line and places it into the appropriate group buckets.
// Lines that are empty or cannot be decoded are silently skipped.
func (g *Grouper) Add(line string) {
	if line == "" {
		return
	}
	var record map[string]any
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for i, rule := range g.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", val)
		if _, exists := g.groups[i][key]; !exists {
			g.groups[i][key] = &group{value: key}
		}
		// Store a copy of the record so later mutations don't affect stored data.
		entry := make(map[string]any, len(record))
		for k, v := range record {
			entry[k] = v
		}
		g.groups[i][key].lines = append(g.groups[i][key].lines, entry)
	}
}

// Flush returns a slice of JSON summary lines — one per (rule, group-value) pair —
// and resets all internal state. The summaries are sorted by field value for
// deterministic output.
func (g *Grouper) Flush() []string {
	g.mu.Lock()
	defer g.mu.Unlock()

	var out []string
	for i, rule := range g.rules {
		outField := rule.OutputField
		if outField == "" {
			outField = rule.Field
		}
		countField := rule.CountField
		if countField == "" {
			countField = "count"
		}

		// Sort keys for deterministic output.
		keys := make([]string, 0, len(g.groups[i]))
		for k := range g.groups[i] {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			grp := g.groups[i][k]
			summary := map[string]any{
				outField:   grp.value,
				countField: len(grp.lines),
			}
			b, err := json.Marshal(summary)
			if err != nil {
				continue
			}
			out = append(out, string(b))
		}
		// Reset bucket for this rule.
		g.groups[i] = make(map[string]*group)
	}
	return out
}
