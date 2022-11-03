package funcs

import (
	"pimtrace"
)

type FunctionDef func(d pimtrace.Entry) (pimtrace.Value, error)

func Functions() map[string]FunctionDef {
	return map[string]FunctionDef{
		"count": Count,
		"sum":   Sum,
		"month": Month,
		"year":  Year,
	}
}
