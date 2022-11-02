package funcs

import "pimtrace/argparsers/basic"

func Functions[T any]() map[string]basic.FunctionDef[T] {
	return map[string]basic.FunctionDef[T]{
		"count": Count[T],
		"sum":   Sum[T],
		"month": Month[T],
		"year":  Year[T],
	}
}
