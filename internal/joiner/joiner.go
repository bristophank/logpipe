// Package joiner merges fields from multiple JSON log lines into a single output line.
package joiner

import (
	"encoding/json"
	"fmt"
)

// Rule defines how to join a field from a secondary line into the primary line.
type Rule struct {
	PrimaryKey   string `json:"primary_key"`
	SecondaryKey string `json:"secondary_key"`
	Fields       []string `json:"fields"`
}

// Joiner holds a lookup table of secondary records keyed by secondary_key value.
type Joiner struct {
	rules  []Rule
	table  map[string]map[string]any
}

// New creates a Joiner with the given rules.
func New(rules []Rule) *Joiner {
	return &Joiner{
		rules: rules,
		table: make(map[string]map[string]any),
	}
}

// Index registers a secondary JSON line into the lookup table for the given rule index.
func (j *Joiner) Index(ruleIdx int, line string) error {
	if ruleIdx < 0 || ruleIdx >= len(j.rules) {
		return fmt.Errorf("joiner: rule index %d out of range", ruleIdx)
	}
	var record map[string]any
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return fmt.Errorf("joiner: invalid JSON: %w", err)
	}
	rule := j.rules[ruleIdx]
	keyVal, ok := record[rule.SecondaryKey]
	if !ok {
		return nil
	}
	keyStr := fmt.Sprintf("%d:%v", ruleIdx, keyVal)
	j.table[keyStr] = record
	return nil
}

// Apply merges indexed secondary fields into the primary JSON line.
func (j *Joiner) Apply(line string) (string, error) {
	if len(j.rules) == 0 {
		return line, nil
	}
	var record map[string]any
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line, nil
	}
	for idx, rule := range j.rules {
		pkVal, ok := record[rule.PrimaryKey]
		if !ok {
			continue
		}
		keyStr := fmt.Sprintf("%d:%v", idx, pkVal)
		secondary, found := j.table[keyStr]
		if !found {
			continue
		}
		for _, f := range rule.Fields {
			if v, exists := secondary[f]; exists {
				record[f] = v
			}
		}
	}
	out, err := json.Marshal(record)
	if err != nil {
		return line, nil
	}
	return string(out), nil
}

// Reset clears the lookup table.
func (j *Joiner) Reset() {
	j.table = make(map[string]map[string]any)
}
