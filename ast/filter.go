package ast

import (
	"pimtrace"

	"github.com/arran4/go-evaluator"
)

func Filter(d pimtrace.Data, expression *evaluator.Query) (pimtrace.Data, error) {
	i, o := 0, 0
	for i+o < d.Len() {
		e := d.Entry(i + o)
		keep := expression.Evaluate(e)
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
