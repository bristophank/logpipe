package sorter

import (
	"encoding/json"
	"sort"
	"strings"
)

// Rule defines a sort operation on a JSON array field.
type Rule struct {
	Field string // field containing the array of objects to sort
	By    string // key within each object to sort by
	Order string // "asc" or "desc" (default: "asc")
}

// Sorter sorts JSON array fields within log lines.
type Sorter struct {
	rules []Rule
}

// New creates a Sorter with the given rules.
func New(rules []Rule) *Sorter {
	return &Sorter{rules: rules}
}

// Apply processes a single JSON log line, sorting array fields per rules.
// Returns the original line unchanged if there are no rules or on parse error.
func (s *Sorter) Apply(line string) string {
	if len(s.rules) == 0 || strings.TrimSpace(line) == "" {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, r := range s.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		arr, ok := val.([]interface{})
		if !ok {
			continue
		}
		sorted := make([]interface{}, len(arr))
		copy(sorted, arr)
		sort.SliceStable(sorted, func(i, j int) bool {
			a := extractStr(sorted[i], r.By)
			b := extractStr(sorted[j], r.By)
			if strings.EqualFold(r.Order, "desc") {
				return a > b
			}
			return a < b
		})
		obj[r.Field] = sorted
	}

	return toString(obj)
}

func extractStr(v interface{}, key string) string {
	m, ok := v.(map[string]interface{})
	if !ok {
		return ""
	}
	if key == "" {
		return ""
	}
	if val, ok := m[key]; ok {
		switch t := val.(type) {
		case string:
			return t
		default:
			b, _ := json.Marshal(t)
			return string(b)
		}
	}
	return ""
}

func toString(obj map[string]interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}
