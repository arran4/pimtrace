package funcs

import (
	"github.com/arran4/go-evaluator"
	"pimtrace"
)

type Hour[T ValueExpression] struct{}

var _ Function[ValueExpression] = Hour[ValueExpression]{}

func (c Hour[T]) Name() string {
	return "hour"
}

func (c Hour[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns the hour of the day (0-23) of that date",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns the hour of the day (0-23) of that date",
		},
	}
}

func (c Hour[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	t, err := temporalCoerce("hour", d, args, ctx)
	if err != nil {
		return &pimtrace.SimpleNilValue{}, err
	}
	if t == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	return pimtrace.SimpleIntegerValue(int(t.Hour())), nil
}
