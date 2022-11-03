package ast

import (
	"pimtrace"
)

func Filter(d pimtrace.Data, expression BooleanExpression) (pimtrace.Data, error) {
	i, o := 0, 0
	for i+o < d.Len() {
		e := d.Entry(i + o)
		keep, err := expression.Execute(e)
		if err != nil {
			return nil, err
		}
		if o > 0 {
			d.SetEntry(i, e)
		}
		if !keep {
			o++
		} else {
			i++
		}
	}
	if o > 0 {
		d = d.Truncate(i)
	}
	return d, nil
}
