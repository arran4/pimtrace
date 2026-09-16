package funcs

import (
	"github.com/arran4/go-evaluator"
	"pimtrace"
)

type Date[T ValueExpression] struct{}

var _ Function[ValueExpression] = Date[ValueExpression]{}

func (c Date[T]) Name() string {
	return "date"
}

func (c Date[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns a canonical date string (YYYY-MM-DD) keeping local timezone",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns a canonical date string (YYYY-MM-DD) keeping local timezone",
		},
	}
}

func (c Date[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	t, err := temporalCoerce("date", d, args, ctx)
	if err != nil {
		return &pimtrace.SimpleNilValue{}, err
	}
	if t == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	// Return a canonical date string representation (e.g. YYYY-MM-DD) keeping local timezone.
	return pimtrace.SimpleStringValue(t.Format("2006-01-02")), nil
}
