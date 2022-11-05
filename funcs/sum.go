package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/groupdata"
)

type Sum[T ValueExpression] struct{}

var _ Function[ValueExpression] = Sum[ValueExpression]{}

func (c Sum[T]) Name() string {
	return "sum"
}

func (c Sum[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Description: "Returns a sum of lines represented by this",
		},
		{
			Args:        []Argument{Any},
			Description: "Returns the number of truthy elements returned",
		},
	}
}

func (c Sum[T]) Run(d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	if len(args) == 0 {
		return pimtrace.SimpleIntegerValue(1), nil
	}
	dd, ok := d.(*groupdata.Row)
	if !ok {
		return pimtrace.SimpleIntegerValue(1), nil
	}
	value := 0
	for i := 0; i < dd.Contents.Len(); i++ {
		e := dd.Contents.Entry(i)
		r, err := args[0].Execute(e)
		if err != nil {
			return nil, err
		}
		i := r.Integer()
		if i != nil {
			value += *i
		}
	}
	return pimtrace.SimpleIntegerValue(value), nil
}
