package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format defines the output format type.
type Format string

const (
	FormatJSON   Format = "json"
	FormatText   Format = "text"
	FormatPretty Format = "pretty"
)

// Formatter transforms a raw JSON log line into the desired output format.
type Formatter struct {
	format Format
}

// New creates a new Formatter for the given format string.
// Defaults to JSON if the format is unrecognized.
func New(format string) *Formatter {
	f := Format(strings.ToLower(format))
	switch f {
	case FormatText, FormatPretty:
	default:
		f = FormatJSON
	}
	return &Formatter{format: f}
}

// Format transforms the input line according to the configured format.
func (f *Formatter) Format(line string) (string, error) {
	switch f.format {
	case FormatPretty:
		return toPretty(line)
	case FormatText:
		return toText(line)
	default:
		return line, nil
	}
}

func toPretty(line string) (string, error) {
	var buf []byte
	var err error
	var m map[string]json.RawMessage
	if err = json.Unmarshal([]byte(line), &m); err != nil {
		return line, fmt.Errorf("pretty: invalid json: %w", err)
	}
	buf, err = json.MarshalIndent(m, "", "  ")
	if err != nil {
		return line, err
	}
	return string(buf), nil
}

func toText(line string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line, fmt.Errorf("text: invalid json: %w", err)
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, m[k]))
	}
	return strings.Join(parts, " "), nil
}
