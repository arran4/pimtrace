package funcs

import (
	"errors"
	"pimtrace"
)

var (
	ErrExpecting1ArgumentOfTypeStringIntOrDate = errors.New("expecting 1 argument, of type string, int or date")
	ErrNumberError                             = errors.New("unknown number")
	ErrUnsupportedType                         = errors.New("unsupported type")
)

func Month[T ValueExpression](d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	t, err := Arg1OnlyToTime(d, args)
	if err != nil {
		return nil, err
	}
	return pimtrace.SimpleIntegerValue(int(t.Month())), nil
}
