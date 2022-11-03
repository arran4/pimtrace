package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/groupdata"
)

func Count[T ValueExpression](d pimtrace.Entry, args []T) (pimtrace.Value, error) {
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
		r, err := args[0].Execute(e)
		if err != nil {
			return nil, err
		}
		if r.Truthy() {
			value += 1
		}
	}
	return pimtrace.SimpleIntegerValue(value), nil
}
