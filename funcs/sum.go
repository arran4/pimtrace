package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/groupdata"
)

func Sum[T ValueExpression](d pimtrace.Entry, args []T) (pimtrace.Value, error) {
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
