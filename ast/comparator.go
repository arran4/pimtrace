package ast

import (
	"fmt"
	"strconv"
	"strings"

	"pimtrace"
)

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
		default:
			return 0, fmt.Errorf("unable to compare against %T", o)
		}
	}

	if otherVal == nil {
		otherVal = &pimtrace.SimpleNilValue{}
	}

	cIsNumeric := isNumeric(c.Value)
	oIsNumeric := isNumeric(otherVal)

	if cIsNumeric && !oIsNumeric {
		on, err := toNumeric(otherVal)
		if err != nil {
			return 0, fmt.Errorf("numeric comparison failed: RHS %q cannot be coerced: %v", otherVal.String(), err)
		}
		return compareNumeric(c.Value, on)
	} else if oIsNumeric && !cIsNumeric {
		cn, err := toNumeric(c.Value)
		if err != nil {
			return 0, fmt.Errorf("numeric comparison failed: LHS %q cannot be coerced: %v", c.Value.String(), err)
		}
		return compareNumeric(cn, otherVal)
	} else if cIsNumeric && oIsNumeric {
		return compareNumeric(c.Value, otherVal)
	}

	ct, cErr := pimtrace.CoerceDate(c.Value)
	ot, oErr := pimtrace.CoerceDate(otherVal)
	if cErr == nil || oErr == nil {
		if cErr != nil {
			return 0, fmt.Errorf("date comparison failed: LHS %q cannot be coerced: %v", c.Value.String(), cErr)
		}
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
	}

	return strings.Compare(c.Value.String(), otherVal.String()), nil
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
