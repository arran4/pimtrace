package ast

import (
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

// ParseDate is a shared, reusable date coercion helper designed to be consumed
// by both the boolean comparison engine and future features (like #53).
func ParseDate(s string) (time.Time, error) {
	// Attempt fast path with common format or delegate to dateparse
	return dateparse.ParseAny(s)
}

// CoercedLiteral represents a literal value that has been typed based on its
// string contents.
type CoercedLiteral struct {
	Original string
	Value    interface{}
}

// ParseLiteral determines the type of a string literal. It checks for integer,
// then float, then date, and falls back to string.
func ParseLiteral(s string) interface{} {
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return int(i)
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if t, err := ParseDate(s); err == nil {
		return t
	}
	return s
}
