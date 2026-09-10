package ast

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"pimtrace"
)

type ComparatorAdapter struct {
	Value pimtrace.Value
}

func (c ComparatorAdapter) Compare(other interface{}) (int, error) {
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
		case float64:
			otherVal = pimtrace.SimpleIntegerValue(int(v))
		default:
			return 0, fmt.Errorf("unable to compare against %T", o)
		}
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

	cIsDate := isDate(c.Value)
	oIsDate := isDate(otherVal)
	if cIsDate || oIsDate {
		ct := c.Value.Time()
		if ct == nil && cIsDate {
			ct = parseDateFallback(c.Value.String())
		} else if ct == nil && !cIsDate {
			ct = parseDateFallback(c.Value.String())
		}

		ot := otherVal.Time()
		if ot == nil && oIsDate {
			ot = parseDateFallback(otherVal.String())
		} else if ot == nil && !oIsDate {
			ot = parseDateFallback(otherVal.String())
		}

		if ct != nil && ot != nil {
			if ct.Before(*ot) {
				return -1, nil
			}
			if ct.After(*ot) {
				return 1, nil
			}
			return 0, nil
		}
	}

	return strings.Compare(c.Value.String(), otherVal.String()), nil
}

func isNumeric(v pimtrace.Value) bool {
    switch v.(type) {
    case pimtrace.SimpleIntegerValue:
        return true
    }

    s := strings.TrimSpace(v.String())
    if s == "" {
        return false
    }

    val := pimtrace.SimpleStringValue(s)
    if val.Float64() != nil {
        return true
    }
    return false
}

func toNumeric(v pimtrace.Value) (pimtrace.Value, error) {
    s := strings.TrimSpace(v.String())
    if s == "" {
        return nil, fmt.Errorf("empty string")
    }

    if isNumeric(v) {
        return pimtrace.SimpleStringValue(s), nil
    }

    val := pimtrace.SimpleStringValue(s)
    if val.Float64() != nil {
        return val, nil
    }
    return nil, fmt.Errorf("not numeric")
}

func compareNumeric(a, b pimtrace.Value) (int, error) {
    sa := strings.TrimSpace(a.String())
    sb := strings.TrimSpace(b.String())
    va := pimtrace.SimpleStringValue(sa)
    vb := pimtrace.SimpleStringValue(sb)

    af := va.Float64()
    bf := vb.Float64()

    if af == nil || bf == nil {
        return 0, fmt.Errorf("failed to get float64 value")
    }

    if *af < *bf {
        return -1, nil
    }
    if *af > *bf {
        return 1, nil
    }
    return 0, nil
}

func isDate(v pimtrace.Value) bool {
    if v.Time() != nil {
        return true
    }
    t := parseDateFallback(v.String())
    return t != nil
}

func parseDateFallback(s string) *time.Time {
    s = strings.TrimSpace(s)
    if s == "" {
        return nil
    }
    t, err := dateparse.ParseAny(s)
    if err == nil {
        return &t
    }
    return nil
}
