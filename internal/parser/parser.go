package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format represents a supported input log format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatLogfmt Format = "logfmt"
	FormatAuto   Format = "auto"
)

// Parser parses raw log lines into key-value maps.
type Parser struct {
	format Format
}

// New creates a Parser for the given format.
func New(format Format) *Parser {
	if format == "" {
		format = FormatAuto
	}
	return &Parser{format: format}
}

// Parse converts a raw line into a map. Returns an error if parsing fails.
func (p *Parser) Parse(line string) (map[string]any, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}
	switch p.format {
	case FormatJSON:
		return parseJSON(line)
	case FormatLogfmt:
		return parseLogfmt(line)
	default: // auto
		if strings.HasPrefix(line, "{") {
			return parseJSON(line)
		}
		return parseLogfmt(line)
	}
}

func parseJSON(line string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return nil, fmt.Errorf("json parse: %w", err)
	}
	return m, nil
}

func parseLogfmt(line string) (map[string]any, error) {
	m := make(map[string]any)
	for _, pair := range strings.Fields(line) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 1 {
			m[parts[0]] = true
			continue
		}
		v := strings.Trim(parts[1], `"`)
		m[parts[0]] = v
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("logfmt parse: no key=value pairs found")
	}
	return m, nil
}
