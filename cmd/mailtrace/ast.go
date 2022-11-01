package main

import (
	"fmt"
	"strings"
)

type Operation interface {
	Execute(d Data) (Data, error)
}

type CompoundStatement struct {
	Statements []Operation
}

func (o *CompoundStatement) Execute(d Data) (Data, error) {
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
	return o
}

var _ Operation = (*CompoundStatement)(nil)

type BooleanExpression interface {
	Execute(d Entry) (bool, error)
}

type NotOp struct {
	Not BooleanExpression
}

func (n *NotOp) Execute(d Entry) (bool, error) {
	v, err := n.Not.Execute(d)
	return !v, err
}

var _ BooleanExpression = (*NotOp)(nil)

type ValueExpression interface {
	Execute(d Entry) (Value, error)
}

type ConstantExpression string

func (ve ConstantExpression) Execute(d Entry) (Value, error) {
	return SimpleStringValue(ve), nil
}

type EntryExpression string

func (ve EntryExpression) Execute(d Entry) (Value, error) {
	return d.Get(string(ve)), nil
}

type OpFunc func(Value, Value) (bool, error)

func EqualOp(rhsv Value, lhsv Value) (bool, error) {
	return rhsv.String() == lhsv.String(), nil
}

var _ OpFunc = EqualOp

func ContainsOp(rhsv Value, lhsv Value) (bool, error) {
	return strings.Contains(rhsv.String(), lhsv.String()), nil
}

var _ OpFunc = ContainsOp

func IContainsOp(rhsv Value, lhsv Value) (bool, error) {
	return strings.Contains(strings.ToLower(rhsv.String()), strings.ToLower(lhsv.String())), nil
}

var _ OpFunc = IContainsOp

type Op struct {
	Op  OpFunc
	LHS ValueExpression
	RHS ValueExpression
}

func (e *Op) Execute(d Entry) (bool, error) {
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

type FilterStatement struct {
	Expression BooleanExpression
}

func (f FilterStatement) Execute(d Data) (Data, error) {
	return Filter(d, f.Expression)
}

var _ Operation = (*FilterStatement)(nil)