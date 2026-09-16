package funcs

import (
	"fmt"
	"pimtrace"
	"time"

	"github.com/arran4/go-evaluator"
)

// temporalCoerce provides a shared date/time coercion policy for temporal functions.
// It accepts one argument, handles PIMTrace values and string values using pimtrace.CoerceDate,
// preserves Unix integer timestamps, and returns useful errors.
func temporalCoerce[T ValueExpression](funcName string, d pimtrace.Entry, args []T, ctx *evaluator.Context) (*time.Time, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("%s: %w", funcName, ErrExpecting1ArgumentOfTypeStringIntOrDate)
	}

	v, err := args[0].Execute(d, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	if v == nil {
		return nil, fmt.Errorf("%s: %w", funcName, ErrEmptyType)
	}

	if _, isNilValue := v.(*pimtrace.SimpleNilValue); isNilValue {
		// Nil/missing values yield nil time without error, propagating the nil semantics.
		return nil, nil
	}

	var t *time.Time
	switch v.(type) {
	case pimtrace.SimpleIntegerValue:
		// Preserve existing Unix integer timestamps behavior.
		i := v.Integer()
		if i == nil {
			return nil, fmt.Errorf("%s parse: %w", funcName, ErrNumberError)
		}
		vt := time.Unix(int64(*i), 0)
		t = &vt
	case pimtrace.SimpleFloatValue:
		// Optionally support floats as unix time, but stick to Integer for now or reject.
		return nil, fmt.Errorf("%s: %w: %s", funcName, ErrUnsupportedType, v.Type())
	default:
		// Delegate everything else (string, date/time) to pimtrace.CoerceDate.
		t, err = pimtrace.CoerceDate(v)
		if err != nil {
			return nil, fmt.Errorf("%s coercion: %w", funcName, err)
		}
	}

	if t == nil {
		return nil, nil
	}

	return t, nil
}
