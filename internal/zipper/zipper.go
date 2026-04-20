// Package zipper merges fields from multiple JSON log lines into a single
// output line, matched by a common key field.
package zipper

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Rule defines how two log lines should be merged.
type Rule struct {
	// Key is the field name used to match lines from Left and Right streams.
	Key string `json:"key"`
	// Prefix is an optional prefix applied to fields merged from the right side.
	Prefix string `json:"prefix"`
}

// Zipper holds buffered entries from a primary stream and merges them with
// entries from a secondary stream when a matching key is found.
type Zipper struct {
	rules []Rule
	mu    sync.Mutex
	// left holds buffered primary-stream entries keyed by rule key value.
	left map[string]map[string]map[string]interface{}
}

// New creates a Zipper with the provided merge rules.
func New(rules []Rule) *Zipper {
	return &Zipper{
		rules: rules,
		left:  make(map[string]map[string]map[string]interface{}),
	}
}

// Buffer stores a primary-stream JSON line so it can later be merged when a
// matching secondary line arrives. Lines that cannot be decoded are silently
// dropped.
func (z *Zipper) Buffer(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return
	}

	z.mu.Lock()
	defer z.mu.Unlock()

	for _, r := range z.rules {
		val, ok := obj[r.Key]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", val)
		if z.left[r.Key] == nil {
			z.left[r.Key] = make(map[string]map[string]interface{})
		}
		z.left[r.Key][key] = obj
	}
}

// Merge attempts to merge a secondary-stream JSON line with a buffered primary
// entry that shares the same key value. If a match is found the merged object
// is returned as a JSON string and the buffered entry is consumed. If no match
// is found the secondary line is returned unchanged.
func (z *Zipper) Merge(line string) (string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	z.mu.Lock()
	defer z.mu.Unlock()

	for _, r := range z.rules {
		val, ok := obj[r.Key]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", val)
		buf, found := z.left[r.Key][key]
		if !found {
			continue
		}
		// Merge: start with buffered (left) fields, then overlay right fields.
		merged := make(map[string]interface{}, len(buf)+len(obj))
		for k, v := range buf {
			merged[k] = v
		}
		for k, v := range obj {
			destKey := k
			if r.Prefix != "" && k != r.Key {
				destKey = r.Prefix + k
			}
			merged[destKey] = v
		}
		delete(z.left[r.Key], key)
		b, err := json.Marshal(merged)
		if err != nil {
			return line, err
		}
		return string(b), nil
	}

	// No match found — return the secondary line as-is.
	return line, nil
}

// Len returns the total number of buffered primary-stream entries.
func (z *Zipper) Len() int {
	z.mu.Lock()
	defer z.mu.Unlock()
	n := 0
	for _, m := range z.left {
		n += len(m)
	}
	return n
}
