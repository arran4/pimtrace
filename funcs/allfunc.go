package funcs

import (
	"pimtrace/ast"
)

func Functions[T any]() map[string]ast.FunctionDef[T] {
	return map[string]ast.FunctionDef[T]{
		"count": Count[T],
		"sum":   Sum[T],
		"month": Month[T],
		"year":  Year[T],
	}
}
