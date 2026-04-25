package stamper

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a field to stamp with a sequential or fixed counter value.
type Rule struct {
	Field  string // target field name
	Start  int    // initial value (default 0)
	Step   int    // increment per line (default 1)
	Format string // "int" or "string" (default "int")
}

// Stamper injects a per-rule counter value into each JSON log line.
type Stamper struct {
	rules    []Rule
	counters []int
}

// New creates a Stamper with the given rules.
// Rules with Step == 0 default to Step = 1.
func New(rules []Rule) *Stamper {
	counters := make([]int, len(rules))
	for i, r := range rules {
		counters[i] = r.Start
		if rules[i].Step == 0 {
			rules[i].Step = 1
		}
		if rules[i].Format == "" {
			rules[i].Format = "int"
		}
	}
	return &Stamper{rules: rules, counters: counters}
}

// Apply stamps each rule's current counter value into the JSON line,
// then advances the counter. Returns the original line on parse error.
func (s *Stamper) Apply(line string) string {
	if len(s.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for i, r := range s.rules {
		val := s.counters[i]
		if r.Format == "string" {
			obj[r.Field] = strconv.Itoa(val)
		} else {
			obj[r.Field] = val
		}
		s.counters[i] += r.Step
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return fmt.Sprintf("%s", out)
}

// Reset resets all counters to their configured Start values.
func (s *Stamper) Reset() {
	for i, r := range s.rules {
		s.counters[i] = r.Start
	}
}
