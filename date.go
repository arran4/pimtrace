package pimtrace

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// ParseDate parses a string representation of a date/time using the project's dateparse policy.
// It trims whitespace and returns an error if the input is empty or cannot be parsed as a date.
func ParseDate(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty date string")
	}
	t, err := dateparse.ParseAny(s)
	if err == nil {
		return &t, nil
	}
	for _, layout := range []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	} {
		if pt, err2 := time.Parse(layout, s); err2 == nil {
			return &pt, nil
		}
	}
	return nil, fmt.Errorf("parse date %q: %w", s, err)
}

// CoerceDate attempts to extract or parse a time.Time from a Value.
// It returns an error if the value is nil, numeric, or cannot be parsed as a date.
func CoerceDate(v Value) (*time.Time, error) {
	if v == nil {
		return nil, fmt.Errorf("empty/nil value")
	}
	switch v.(type) {
	case SimpleIntegerValue, SimpleFloatValue:
		return nil, fmt.Errorf("value %s is numeric, not a date", v.String())
	}
	if t := v.Time(); t != nil {
		return t, nil
	}
	return ParseDate(v.String())
}
