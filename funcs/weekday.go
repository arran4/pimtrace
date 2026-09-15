package funcs

import (
	"github.com/arran4/go-evaluator"
	"pimtrace"
)

type Weekday[T ValueExpression] struct{}

var _ Function[ValueExpression] = Weekday[ValueExpression]{}

func (c Weekday[T]) Name() string {
	return "weekday"
}

func (c Weekday[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns the weekday name (e.g. Monday) of that date",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns the weekday name (e.g. Monday) of that date",
		},
	}
}

func (c Weekday[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	t, err := temporalCoerce("weekday", d, args, ctx)
	if err != nil {
		return &pimtrace.SimpleNilValue{}, err
	}
	if t == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	return pimtrace.SimpleStringValue(t.Weekday().String()), nil
}
