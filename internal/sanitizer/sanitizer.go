package sanitizer

import (
	"encoding/json"
	"strings"
)

// Rule defines a sanitization rule for a specific field.
type Rule struct {
	Field   string // JSON field to sanitize
	Mode    string // "strip_html", "trim", "collapse_spaces", "alphanumeric"
	Fallback string // optional: value to use if field is empty after sanitizing
}

// Sanitizer cleans string fields in JSON log lines.
type Sanitizer struct {
	rules []Rule
}

// New creates a Sanitizer with the given rules.
func New(rules []Rule) *Sanitizer {
	return &Sanitizer{rules: rules}
}

// Apply sanitizes fields in the JSON line according to configured rules.
// Returns the original line unchanged if it is empty or not valid JSON.
func (s *Sanitizer) Apply(line string) string {
	if len(s.rules) == 0 || strings.TrimSpace(line) == "" {
		return line
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	for _, r := range s.rules {
		v, ok := m[r.Field]
		if !ok {
			continue
		}
		str, ok := v.(string)
		if !ok {
			continue
		}
		str = applyMode(str, r.Mode)
		if str == "" && r.Fallback != "" {
			str = r.Fallback
		}
		m[r.Field] = str
	}
	return toString(m)
}

func applyMode(s, mode string) string {
	switch mode {
	case "strip_html":
		return stripHTML(s)
	case "trim":
		return strings.TrimSpace(s)
	case "collapse_spaces":
		parts := strings.Fields(s)
		return strings.Join(parts, " ")
	case "alphanumeric":
		var b strings.Builder
		for _, ch := range s {
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
				b.WriteRune(ch)
			}
		}
		return b.String()
	default:
		return s
	}
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, ch := range s {
		if ch == '<' {
			inTag = true
			continue
		}
		if ch == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func toString(m map[string]any) string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}
