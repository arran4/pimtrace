package maildata

import (
	"pimtrace"
	"pimtrace/ast"
)

type MBoxOutput struct{}

func (M *MBoxOutput) Execute(d pimtrace.Data) (pimtrace.Data, error) {
	//TODO implement me
	panic("implement me")
}

var _ ast.Operation = (*MBoxOutput)(nil)
