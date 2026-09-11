package ast

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"pimtrace"
)

type comparisonMode int

const (
	modeLexical comparisonMode = iota
	modeNumeric
	modeDate
)

func determineComparisonMode(a, b pimtrace.Value) comparisonMode {
	isNumA := isNumeric(a)
	isNumB := isNumeric(b)

	_, dateErrA := pimtrace.CoerceDate(a)
	isDateA := dateErrA == nil
	_, dateErrB := pimtrace.CoerceDate(b)
	isDateB := dateErrB == nil

	// If both operands are date-coercible and at least one operand provides a non-numeric date signal,
	// use DATE comparison.
	// E.g. "20200102" vs "2020-01-03", "20200102" vs RFC3339 date, "2020-01-01" vs "2020-01-03".
	if isDateA && isDateB && (!isNumA || !isNumB) {
		return modeDate
	}

	// Otherwise, if either operand establishes numeric mode, require both operands to coerce
	// numerically and compare numerically.
	if isNumA || isNumB {
		return modeNumeric
	}

	// Otherwise, if either operand establishes date mode, require both operands to coerce as dates
	// or return a useful date coercion error.
	if isDateA || isDateB {
		return modeDate
	}

	// Only use lexical string ordering when neither typed mode is established.
	return modeLexical
}

type ComparatorAdapter struct {
	Value pimtrace.Value
}

func (c ComparatorAdapter) Compare(other interface{}) (int, error) {
	if c.Value == nil {
		c.Value = &pimtrace.SimpleNilValue{}
	}

	var otherVal pimtrace.Value
	switch o := other.(type) {
	case ComparatorAdapter:
		otherVal = o.Value
	case pimtrace.Value:
		otherVal = o
	default:
		switch v := o.(type) {
		case string:
			otherVal = pimtrace.SimpleStringValue(v)
		case int:
			otherVal = pimtrace.SimpleIntegerValue(v)
		case int64:
			otherVal = pimtrace.SimpleIntegerValue(int(v))
		case float64:
			otherVal = pimtrace.SimpleFloatValue(v)
		case float32:
			otherVal = pimtrace.SimpleFloatValue(float64(v))
		case time.Time:
			otherVal = pimtrace.SimpleStringValue(v.Format(time.RFC3339Nano))
		case *time.Time:
			if v != nil {
				otherVal = pimtrace.SimpleStringValue(v.Format(time.RFC3339Nano))
			} else {
				otherVal = &pimtrace.SimpleNilValue{}
			}
		default:
			return 0, fmt.Errorf("unable to compare against %T", o)
		}
	}

	if otherVal == nil {
		otherVal = &pimtrace.SimpleNilValue{}
	}

	mode := determineComparisonMode(c.Value, otherVal)
	switch mode {
	case modeDate:
		ct, cErr := pimtrace.CoerceDate(c.Value)
		if cErr != nil {
			return 0, fmt.Errorf("date comparison failed: LHS %q cannot be coerced: %v", c.Value.String(), cErr)
		}
		ot, oErr := pimtrace.CoerceDate(otherVal)
		if oErr != nil {
			return 0, fmt.Errorf("date comparison failed: RHS %q cannot be coerced: %v", otherVal.String(), oErr)
		}
		if ct.Before(*ot) {
			return -1, nil
		}
		if ct.After(*ot) {
			return 1, nil
		}
		return 0, nil

	case modeNumeric:
		cn, err := toNumeric(c.Value)
		if err != nil {
			return 0, fmt.Errorf("numeric comparison failed: LHS %q cannot be coerced: %v", c.Value.String(), err)
		}
		on, err := toNumeric(otherVal)
		if err != nil {
			return 0, fmt.Errorf("numeric comparison failed: RHS %q cannot be coerced: %v", otherVal.String(), err)
		}
		return compareNumeric(cn, on)

	case modeLexical:
		return strings.Compare(c.Value.String(), otherVal.String()), nil

	default:
		return strings.Compare(c.Value.String(), otherVal.String()), nil
	}
}

func getFloat64(v pimtrace.Value) *float64 {
	if v == nil {
		return nil
	}
	switch tv := v.(type) {
	case pimtrace.SimpleFloatValue:
		f := float64(tv)
		return &f
	case pimtrace.SimpleIntegerValue:
		f := float64(tv)
		return &f
	}
	if f := v.Float64(); f != nil {
		return f
	}
	s := strings.TrimSpace(v.String())
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

func isNumeric(v pimtrace.Value) bool {
	return getFloat64(v) != nil
}

func toNumeric(v pimtrace.Value) (pimtrace.Value, error) {
	if v == nil {
		return nil, fmt.Errorf("empty/nil value")
	}
	switch tv := v.(type) {
	case pimtrace.SimpleFloatValue, pimtrace.SimpleIntegerValue:
		return tv, nil
	}
	s := strings.TrimSpace(v.String())
	if s == "" {
		return nil, fmt.Errorf("empty string")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("not numeric: %w", err)
	}
	return pimtrace.SimpleFloatValue(f), nil
}

func compareNumeric(a, b pimtrace.Value) (int, error) {
	af := getFloat64(a)
	bf := getFloat64(b)

	if af == nil || bf == nil {
		return 0, fmt.Errorf("failed to get numeric value")
	}

	if *af < *bf {
		return -1, nil
	}
	if *af > *bf {
		return 1, nil
	}
	return 0, nil
}
