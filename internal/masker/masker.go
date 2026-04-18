package masker

import (
	"encoding/json"
	"strings"
)

// Rule defines a field and the mask character to apply.
type Rule struct {
	Field string `json:"field"`
	Mask  string `json:"mask"`   // default "***"
	Keep  int    `json:"keep"`   // number of trailing chars to keep
}

// Masker applies masking rules to JSON log lines.
type Masker struct {
	rules []Rule
}

// New creates a Masker with the given rules.
func New(rules []Rule) *Masker {
	return &Masker{rules: rules}
}

// Apply masks fields in the JSON line according to rules.
// Lines that are not valid JSON are returned unchanged.
func (m *Masker) Apply(line string) string {
	if len(m.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range m.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		str, ok := v.(string)
		if !ok {
			continue
		}
		mask := r.Mask
		if mask == "" {
			mask = "***"
		}
		if r.Keep > 0 && r.Keep < len(str) {
			obj[r.Field] = mask + str[len(str)-r.Keep:]
		} else {
			obj[r.Field] = mask
		}
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return strings.TrimSpace(string(b))
}
