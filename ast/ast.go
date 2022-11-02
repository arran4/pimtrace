package ast

import (
	"fmt"
	"pimtrace"
	"pimtrace/util"
	"strings"
	"unicode"
)

type Operation interface {
	Execute(d pimtrace.Data[T]) (pimtrace.Data[T], error)
}

type CompoundStatement struct {
	Statements []Operation
}

func (o *CompoundStatement) Execute(d pimtrace.Data[T]) (pimtrace.Data[T], error) {
	for _, op := range o.Statements {
		var err error
		d, err = op.Execute(d)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (o *CompoundStatement) Simplify() Operation {
	if len(o.Statements) == 0 {
		return nil
	}
	if len(o.Statements) == 1 {
		return o.Statements[0]
	}
	var result []Operation
	for i, statement := range o.Statements {
		switch statement := statement.(type) {
		case *CompoundStatement:
			if result == nil {
				result = append([]Operation{}, o.Statements[:i]...)
			}
			result = append(result, statement.Statements...)
		default:
			if result == nil {
				continue
			}
			result = append(result, statement)
		}
	}
	if result != nil {
		o.Statements = result
	}
	return o
}

var _ Operation = (*CompoundStatement)(nil)

type BooleanExpression[T any] interface {
	Execute(d pimtrace.Entry[T]) (bool, error)
}

type NotOp[T any] struct {
	Not BooleanExpression[T]
}

func (n *NotOp[T]) Execute(d pimtrace.Entry[T]) (bool, error) {
	v, err := n.Not.Execute(d)
	return !v, err
}

var _ BooleanExpression[any] = (*NotOp[any])(nil)

type ValueExpression[T any] interface {
	Execute(d pimtrace.Entry[T]) (pimtrace.Value, error)
	ColumnName() string
}

type ConstantExpression[T any] string

func (ve ConstantExpression[T]) ColumnName() string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return '-'
	}, string(ve))
}

func (ve ConstantExpression[T]) Execute(d pimtrace.Entry[T]) (pimtrace.Value, error) {
	return pimtrace.SimpleStringValue(ve), nil
}

type EntryExpression[T any] string

func (ve EntryExpression[T]) ColumnName() string {
	ss := strings.SplitN(string(ve), ".", 2)
	s := ""
	if len(ss) > 1 {
		s = ss[1]
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return '-'
	}, s)
}

func (ve EntryExpression[T]) Execute(d pimtrace.Entry[T]) (pimtrace.Value, error) {
	return d.Get(string(ve)), nil
}

type OpFunc func(pimtrace.Value, pimtrace.Value) (bool, error)

func EqualOp(rhsv pimtrace.Value, lhsv pimtrace.Value) (bool, error) {
	return rhsv.String() == lhsv.String(), nil
}

var _ OpFunc = EqualOp

func ContainsOp(rhsv pimtrace.Value, lhsv pimtrace.Value) (bool, error) {
	return strings.Contains(rhsv.String(), lhsv.String()), nil
}

var _ OpFunc = ContainsOp

func IContainsOp(rhsv pimtrace.Value, lhsv pimtrace.Value) (bool, error) {
	return strings.Contains(strings.ToLower(rhsv.String()), strings.ToLower(lhsv.String())), nil
}

var _ OpFunc = IContainsOp

type Op[T any] struct {
	Op  OpFunc
	LHS ValueExpression[T]
	RHS ValueExpression[T]
}

func (e *Op[T]) Execute(d pimtrace.Entry[T]) (bool, error) {
	if e.LHS == nil {
		return false, fmt.Errorf("LHS invalid issue with Op")
	}
	if e.RHS == nil {
		return false, fmt.Errorf("RHS invalid with Op")
	}
	if e.Op == nil {
		return false, fmt.Errorf("op invalid with Op")
	}
	lhsv, err := e.LHS.Execute(d)
	if err != nil {
		return false, fmt.Errorf("LHS error: %w", err)
	}
	rhsv, err := e.RHS.Execute(d)
	if err != nil {
		return false, fmt.Errorf("RHS error: %w", err)
	}
	return e.Op(rhsv, lhsv)
}

type FilterStatement[T any] struct {
	Expression BooleanExpression[T]
}

func (f FilterStatement[T]) Execute(d pimtrace.Data[T]) (pimtrace.Data[T], error) {
	return util.Filter(d, f.Expression)
}

var _ Operation = (*FilterStatement)(nil)

type MBoxOutput struct{}

func (M *MBoxOutput) Execute(d pimtrace.Data[T]) (pimtrace.Data[T], error) {
	//TODO implement me
	panic("implement me")
}

var _ Operation = (*MBoxOutput)(nil)
