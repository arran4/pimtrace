package funcs

import (
	"github.com/arran4/go-evaluator"
	"pimtrace"
)

type Day[T ValueExpression] struct{}

var _ Function[ValueExpression] = Day[ValueExpression]{}

func (c Day[T]) Name() string {
	return "day"
}

func (c Day[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns the day of the month of that date",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns the day of the month of that date",
		},
	}
}

func (c Day[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	t, err := temporalCoerce("day", d, args, ctx)
	if err != nil {
		return &pimtrace.SimpleNilValue{}, err
	}
	if t == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	return pimtrace.SimpleIntegerValue(int(t.Day())), nil
}
