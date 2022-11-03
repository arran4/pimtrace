package funcs

import (
	"pimtrace"
)

type FunctionDef[T any] func(d pimtrace.Entry[T]) (pimtrace.Value, error)

func Functions[T any]() map[string]FunctionDef[T] {
	return map[string]FunctionDef[T]{
		"count": Count[T],
		"sum":   Sum[T],
		"month": Month[T],
		"year":  Year[T],
	}
}
