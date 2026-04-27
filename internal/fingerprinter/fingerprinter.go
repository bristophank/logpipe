// Package fingerprinter attaches a short deterministic hash to each log line
// based on a configurable set of fields. This makes downstream dedup or
// correlation easier without shipping full payloads.
package fingerprinter

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Rule describes which fields to include in the fingerprint and where to write
// the result.
type Rule struct {
	// Fields is the ordered list of field names whose values form the fingerprint.
	// If empty, all fields are used (sorted by key).
	Fields []string `json:"fields"`
	// Target is the output field name. Defaults to "_fp".
	Target string `json:"target"`
	// Len is the number of hex characters to keep (max 64). Defaults to 16.
	Len int `json:"len"`
}

// Fingerprinter adds fingerprint fields to JSON log lines.
type Fingerprinter struct {
	rules []Rule
}

// New creates a Fingerprinter with the given rules.
func New(rules []Rule) *Fingerprinter {
	return &Fingerprinter{rules: rules}
}

// Apply processes a single JSON line and returns it with fingerprint fields
// added. Non-JSON or empty lines are returned unchanged.
func (f *Fingerprinter) Apply(line string) string {
	if len(f.rules) == 0 {
		return line
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, r := range f.rules {
		target := r.Target
		if target == "" {
			target = "_fp"
		}
		keepLen := r.Len
		if keepLen <= 0 || keepLen > 64 {
			keepLen = 16
		}
		obj[target] = fingerprint(obj, r.Fields, keepLen)
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// fingerprint builds a hex digest from the selected fields of obj.
func fingerprint(obj map[string]interface{}, fields []string, length int) string {
	if len(fields) == 0 {
		fields = make([]string, 0, len(obj))
		for k := range obj {
			fields = append(fields, k)
		}
		sort.Strings(fields)
	}

	var sb strings.Builder
	for _, k := range fields {
		v, ok := obj[k]
		if !ok {
			continue
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(fmt.Sprintf("%v", v))
		sb.WriteByte(';')
	}

	sum := sha256.Sum256([]byte(sb.String()))
	hex := fmt.Sprintf("%x", sum[:])
	if length > len(hex) {
		length = len(hex)
	}
	return hex[:length]
}
