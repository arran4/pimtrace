package funcs

import (
	"errors"
	"pimtrace"
)

var (
	ErrExpecting1ArgumentOfTypeStringIntOrDate = errors.New("expecting 1 argument, of type string, int or date")
	ErrNumberError                             = errors.New("unknown number")
	ErrUnsupportedType                         = errors.New("unsupported type")
	ErrEmptyType                               = errors.New("empty type")
)

type Month[T ValueExpression] struct{}

var _ Function[ValueExpression] = Month[ValueExpression]{}

func (c Month[T]) Name() string {
	return "month"
}

func (c Month[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns the month number of that date",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns the month number of that date",
		},
	}
}

func (c Month[T]) Run(d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	t, err := Arg1OnlyToTime("month", d, args)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return pimtrace.Nil, nil
	}
	return pimtrace.SimpleIntegerValue(int(t.Month())), nil
}
