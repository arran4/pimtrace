package funcs

import (
	"errors"
	"fmt"
	"pimtrace"
	"pimtrace/dataformats/nildata"
)

var (
	ErrExpecting2ArgumentsAnyString = errors.New("expecting 2 arguments: the value => any, the name => string")
)

type As[T ValueExpression] struct{}

var _ Function[ValueExpression] = As[ValueExpression]{}
var _ ColumnNamer[ValueExpression] = As[ValueExpression]{}

func (c As[T]) Name() string {
	return "as"
}

func (c As[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{Any, String},
			Description: "Renames the column to a specific name",
		},
	}
}

func (c As[T]) ColumnName(args []T) string {
	if len(args) < 2 {
		return ""
	}
	v, _ := args[1].Execute(&nildata.Row{})
	return v.String()
}

func (c As[T]) Run(d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("as: %w", ErrExpecting2ArgumentsAnyString)
	}
	return args[0].Execute(d)
}
