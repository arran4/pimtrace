package funcs

import (
	"fmt"
	"github.com/araddon/dateparse"
	"pimtrace"
	"time"
)

func Year[T ValueExpression](d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	t, err := Arg1OnlyToTime("year", d, args)
	if err != nil {
		return nil, err
	}
	return pimtrace.SimpleIntegerValue(int(t.Year())), nil
}

func Arg1OnlyToTime[T ValueExpression](funcName string, d pimtrace.Entry, args []T) (time.Time, error) {
	if len(args) == 0 {
		return time.Time{}, fmt.Errorf("%w", ErrExpecting1ArgumentOfTypeStringIntOrDate)
	}
	v, err := args[0].Execute(d)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", funcName, err)
	}
	if v == nil {
		return time.Time{}, fmt.Errorf("%s: %w", funcName, ErrEmptyType)
	}
	var t time.Time
	switch v.(type) {
	case pimtrace.SimpleIntegerValue:
		i := v.Integer()
		if i == nil {
			return time.Time{}, fmt.Errorf("%s parse: %w", funcName, ErrNumberError)
		}
		t = time.Unix(int64(*i), 0)
	case pimtrace.SimpleStringValue:
		t, err = dateparse.ParseStrict(v.String())
		if err != nil {
			return time.Time{}, fmt.Errorf("%s parse: %w", funcName, err)
		}
	default:
		return time.Time{}, fmt.Errorf("%s %w: %s", funcName, ErrUnsupportedType, v.Type())
	}
	return t, nil
}
