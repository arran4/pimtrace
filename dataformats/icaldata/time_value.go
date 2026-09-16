package icaldata

import (
	"pimtrace"
	"time"
)

type ICalTimeValue struct {
	T              time.Time
	OriginalString string
}

func (s ICalTimeValue) Truthy() bool {
	return !s.T.IsZero()
}

func (s ICalTimeValue) Elements() int {
	return 1
}

func (s ICalTimeValue) Length() int {
	return 1
}

func (s ICalTimeValue) Array() []pimtrace.Value {
	return []pimtrace.Value{s}
}

func (s ICalTimeValue) StringArray() []string {
	return []string{s.String()}
}

func (s ICalTimeValue) Less(jv pimtrace.Value) bool {
	if vt := jv.Time(); vt != nil {
		return s.T.Before(*vt)
	}
	return s.String() < jv.String()
}

func (s ICalTimeValue) Equal(jv pimtrace.Value) bool {
	if vt := jv.Time(); vt != nil {
		return s.T.Equal(*vt)
	}
	return s.String() == jv.String()
}

func (s ICalTimeValue) Time() *time.Time {
	return &s.T
}

func (s ICalTimeValue) Integer() *int {
	// Behave like the original string value for numeric coercions.
	// Only return Integer if OriginalString is actually numeric.
	return pimtrace.SimpleStringValue(s.String()).Integer()
}

func (s ICalTimeValue) Float64() *float64 {
	return pimtrace.SimpleStringValue(s.String()).Float64()
}

func (s ICalTimeValue) Type() pimtrace.Type {
	return pimtrace.String
}

func (s ICalTimeValue) String() string {
	// Preserve the original unparsed property string behavior for users
	// unless they are explicitly testing dates.
	if s.OriginalString != "" {
		return s.OriginalString
	}
	return s.T.Format(time.RFC3339)
}

var _ pimtrace.Value = ICalTimeValue{}
