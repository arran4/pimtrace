package ast

import (
	"fmt"
	"pimtrace"
	"time"

	"github.com/arran4/go-evaluator"
	"github.com/araddon/dateparse"
)

type ComparatorAdapter struct {
	Value pimtrace.Value
}

func (c *ComparatorAdapter) Compare(other interface{}) (int, error) {
	if c.Value == nil {
		return 0, fmt.Errorf("cannot compare nil value")
	}

	switch otherVal := other.(type) {
	case *ComparatorAdapter:
		// When both are PIMTrace values, fallback to pimtrace.Value Less/Equal semantics
		if c.Value.Equal(otherVal.Value) {
			return 0, nil
		}
		if c.Value.Less(otherVal.Value) {
			return -1, nil
		}
		return 1, nil

	case int:
		val := c.Value.Integer()
		if val == nil {
			return 0, fmt.Errorf("could not coerce %v to integer for comparison", c.Value)
		}
		if *val == otherVal {
			return 0, nil
		}
		if *val < otherVal {
			return -1, nil
		}
		return 1, nil

	case int64:
		val := c.Value.Integer()
		if val == nil {
			return 0, fmt.Errorf("could not coerce %v to integer for comparison", c.Value)
		}
		otherInt := int(otherVal)
		if *val == otherInt {
			return 0, nil
		}
		if *val < otherInt {
			return -1, nil
		}
		return 1, nil

	case float64:
		val := c.Value.Float64()
		if val == nil {
			return 0, fmt.Errorf("could not coerce %v to float64 for comparison", c.Value)
		}
		if *val == otherVal {
			return 0, nil
		}
		if *val < otherVal {
			return -1, nil
		}
		return 1, nil

	case float32:
		val := c.Value.Float64()
		if val == nil {
			return 0, fmt.Errorf("could not coerce %v to float64 for comparison", c.Value)
		}
		otherFloat := float64(otherVal)
		if *val == otherFloat {
			return 0, nil
		}
		if *val < otherFloat {
			return -1, nil
		}
		return 1, nil

	case string:
		str := c.Value.String()
		if str == otherVal {
			return 0, nil
		}
		if str < otherVal {
			return -1, nil
		}
		return 1, nil

	case bool:
		truthy := c.Value.Truthy()
		if truthy == otherVal {
			return 0, nil
		}
		if !truthy && otherVal {
			return -1, nil
		}
		return 1, nil

	case time.Time:
		val := c.Value.Time()
		if val == nil {
			// fallback to dateparse
			s := c.Value.String()
			t, err := dateparse.ParseAny(s)
			if err != nil {
				return 0, fmt.Errorf("could not coerce %v to time.Time for comparison: %w", c.Value, err)
			}
			val = &t
		}

		if val.Equal(otherVal) {
			return 0, nil
		}
		if val.Before(otherVal) {
			return -1, nil
		}
		return 1, nil
	}

	return 0, fmt.Errorf("unsupported comparison type %T", other)
}

var _ evaluator.Comparator = (*ComparatorAdapter)(nil)
