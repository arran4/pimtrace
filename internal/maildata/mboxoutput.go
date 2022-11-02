package maildata

import (
	"pimtrace"
	"pimtrace/ast"
)

type MBoxOutput struct{}

func (M *MBoxOutput) Execute(d pimtrace.Data[*MailWithSource]) (pimtrace.Data[*MailWithSource], error) {
	//TODO implement me
	panic("implement me")
}

var _ ast.Operation[*MailWithSource] = (*MBoxOutput)(nil)
