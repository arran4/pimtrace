package funcs

import (
	"pimtrace"
)

type ValueExpression interface {
	Execute(d pimtrace.Entry) (pimtrace.Value, error)
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
	for _, f := range []Function[T]{
		Count[T]{},
		Sum[T]{},
		Month[T]{},
		Year[T]{},
	} {
		m[f.Name()] = f
	}
	return m
}
