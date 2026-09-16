package funcs

import (
	"github.com/arran4/go-evaluator"
	"pimtrace"
)

type Year[T ValueExpression] struct{}

var _ Function[ValueExpression] = Year[ValueExpression]{}

func (c Year[T]) Name() string {
	return "year"
}

func (c Year[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns the year number of that date",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns the year number of that date",
		},
	}
}

func (c Year[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	t, err := temporalCoerce("year", d, args, ctx)
	if err != nil {
		return &pimtrace.SimpleNilValue{}, err
	}
	if t == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	return pimtrace.SimpleIntegerValue(int(t.Year())), nil
}
