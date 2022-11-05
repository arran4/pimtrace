package funcs

import (
	"fmt"
	"github.com/araddon/dateparse"
	"pimtrace"
	"time"
)

func Year[T ValueExpression](d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	t, err := Arg1OnlyToTime(d, args)
	if err != nil {
		return nil, err
	}
	return pimtrace.SimpleIntegerValue(int(t.Year())), nil
}

func Arg1OnlyToTime[T ValueExpression](d pimtrace.Entry, args []T) (time.Time, error) {
	if len(args) == 0 {
		return time.Time{}, fmt.Errorf("%w", ErrExpecting1ArgumentOfTypeStringIntOrDate)
	}
	v, err := args[0].Execute(d)
	if err != nil {
		return time.Time{}, fmt.Errorf("month: %w", err)
	}
	var t time.Time
	switch v.(type) {
	case pimtrace.SimpleIntegerValue:
		i := v.Integer()
		if i == nil {
			return time.Time{}, fmt.Errorf("month parse: %w", ErrNumberError)
		}
		t = time.Unix(int64(*i), 0)
	case pimtrace.SimpleStringValue:
		t, err = dateparse.ParseStrict(v.String())
		if err != nil {
			return time.Time{}, fmt.Errorf("month parse: %w", err)
		}
	default:
		return time.Time{}, fmt.Errorf("%w: %s", ErrUnsupportedType, v.Type())
	}
	return t, nil
}
