package funcs

import (
	"pimtrace"
)

type ValueExpression interface {
	Execute(d pimtrace.Entry) (pimtrace.Value, error)
}

type FunctionDef[T ValueExpression] func(d pimtrace.Entry, args []T) (pimtrace.Value, error)

func Functions[T ValueExpression]() map[string]FunctionDef[T] {
	return map[string]FunctionDef[T]{
		"count": Count[T],
		"sum":   Sum[T],
		"month": Month[T],
		"year":  Year[T],
	}
}
