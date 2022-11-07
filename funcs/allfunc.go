package funcs

import (
	"pimtrace"
)

type ValueExpression interface {
	Execute(d pimtrace.Entry) (pimtrace.Value, error)
}

type ColumnNamer[T ValueExpression] interface {
	ColumnName(args []T) string
}

type FunctionDef[T ValueExpression] func(d pimtrace.Entry, args []T) (pimtrace.Value, error)

type Argument int

func (a Argument) String() string {
	switch a {
	case String:
		return "String"
	case Integer:
		return "Integer"
	case Array:
		return "Array"
	case Any:
		return "Any"
	}
	return "unknown"
}

const (
	String Argument = iota
	Integer
	Array
	Any
)

type ArgumentList struct {
	Args        []Argument
	Description string
}

type Function[T ValueExpression] interface {
	Name() string
	Arguments() []ArgumentList
	Run(d pimtrace.Entry, args []T) (pimtrace.Value, error)
}

func Functions[T ValueExpression]() map[string]Function[T] {
	m := map[string]Function[T]{}
	// Don't forget to run: cmd/docs/genfunctionmd every time you update anything in this package.
	for _, f := range []Function[T]{
		Count[T]{},
		Sum[T]{},
		Month[T]{},
		Year[T]{},
		As[T]{},
	} {
		m[f.Name()] = f
	}
	return m
}
