package pimtrace

import (
	"strings"
)

type FieldDescriptor struct {
	Name            string
	QueryExpression string
	Component       string
}

type SchemaDescriber interface {
	SchemaFields() []FieldDescriptor
}

// IsQueryable checks if a field name can be safely used in the current basic parser
// without introducing syntax errors or unintended token splitting.
func IsQueryable(s string) bool {
	if s == "" {
		return false
	}
	if strings.ContainsAny(s, " \t\n\r()\"'=<>!&|[]") {
		return false
	}
	return true
}
