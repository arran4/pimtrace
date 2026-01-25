package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/groupdata"

	"github.com/arran4/go-evaluator"
)

type Count[T ValueExpression] struct{}

var _ Function[ValueExpression] = Count[ValueExpression]{}

func (c Count[T]) Name() string {
	return "count"
}

func (c Count[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Description: "Returns a count of lines represented by this",
		},
		{
			Args:        []Argument{Any},
			Description: "Returns the number of truthy elements returned",
		},
	}
}

func (c Count[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	dd, ok := d.(*groupdata.Row)
	if !ok {
		return pimtrace.SimpleIntegerValue(1), nil
	}
	if len(args) == 0 {
		return pimtrace.SimpleIntegerValue(dd.Contents.Len()), nil
	}
	value := 0
	for i := 0; i < dd.Contents.Len(); i++ {
		e := dd.Contents.Entry(i)
		r, err := args[0].Execute(e, ctx)
		if err != nil {
			return nil, err
		}
		if r.Truthy() {
			value += 1
		}
	}
	return pimtrace.SimpleIntegerValue(value), nil
}
