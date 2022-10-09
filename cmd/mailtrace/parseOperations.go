package main

import (
	"fmt"
	"reflect"
	"strings"
)

var ErrParserNothingFound = fmt.Errorf("no token found")

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

var _ BooleanExpression = (*NotOp)(nil)

type ValueExpression interface {
	Execute(d Entry) (Value, error)
}

type ConstantExpression string

func (ve ConstantExpression) Execute(d Entry) (Value, error) {
	return SimpleStringValue(ve), nil
}

func (n *NotOp) Execute(d Entry) (bool, error) {
	v, err := n.Not.Execute(d)
	return !v, err
}

type FilterEquals string
type EntryExpression string
type FilterNot string
type FilterTerminator string

func FilterIdentify(s string) (any, error) {
	ss := strings.SplitN(s, ".", 2)
	switch ss[0] {
	case "map", "filter", "where":
		return FilterTerminator(s), nil
	case "not":
		return FilterNot(s), nil
	case "eq":
		return FilterEquals(s), nil
	case "h", "header":
		return EntryExpression(s), nil
	case "":
		if strings.HasPrefix(s, ".") {
			return ConstantExpression(ss[1]), nil
		}
	}
	return nil, fmt.Errorf("unknown token: %s", ss[0])
}

func FilterTokenizerScanN(args []string, n int) ([]any, []string, error) {
	i := 0
	r := []any{}
	for ; i < n && i < len(args); i++ {
		t, err := FilterIdentify(args[i])
		if err != nil {
			return nil, nil, err
		}
		r = append(r, t)
	}
	return r, args[i:], nil
}

type EqualOp struct {
	LHS ValueExpression
	RHS ValueExpression
}

func (e *EqualOp) Execute(d Entry) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func ParseFilter(args []string, statements []Operation) (BooleanExpression, []string, error) {
	p := args
	tks, remain, err := FilterTokenizerScanN(args, 3)
	if err != nil {
		return nil, nil, err
	}
	if FilterTokenMatcher(tks, FilterNot("")) {
		var op BooleanExpression
		op, remain, err = ParseFilter(args[1:], []Operation{})
		if err != nil {
			return nil, nil, err
		}
		p = remain
		return &NotOp{
			Not: op,
		}, p, nil
	}
	if FilterTokenMatcher(tks, []any{EntryExpression(""), ConstantExpression("")}, FilterEquals(""), []any{EntryExpression(""), ConstantExpression("")}) {
		p = remain
		return &EqualOp{
			LHS: tks[0].(ValueExpression),
			RHS: tks[2].(ValueExpression),
		}, p, nil
	}

	return nil, nil, fmt.Errorf("at %s: %w", p[0], ErrParserNothingFound)
}

func FilterTokenMatcher(tks []any, tokenTypes ...any) bool {
	for i := 0; i < len(tokenTypes); i++ {
		if i >= len(tks) && i < len(tokenTypes) {
			return false
		}
		switch reflect.TypeOf(tokenTypes[i]).Kind() {
		case reflect.Slice:
			for _, stt := range (tokenTypes[i]).([]any) {
				if reflect.TypeOf(tks[i]) != reflect.TypeOf(stt) {
					return false
				}
			}
		case reflect.String:
			if reflect.TypeOf(tks[i]) != reflect.TypeOf(tokenTypes[i]) {
				return false
			}
		}
	}
	return true
}

type FilterStatement struct {
	Expression BooleanExpression
}

func (f FilterStatement) Execute(d Data) (Data, error) {
	return Filter(d, f.Expression)
}

var _ Operation = (*FilterStatement)(nil)

func ParseFilters(args []string) (Operation, []string, error) {
	result := &CompoundStatement{}
	p := args
	for len(p) > 0 {
		switch p[0] {
		case "map":
			return result.Simplify(), p, nil
		case "filter", "where":
			fallthrough
		default:
			boolExp, remain, err := ParseFilter(p[1:], result.Statements)
			if err != nil {
				return nil, nil, err
			}
			p = remain
			statement := &FilterStatement{
				Expression: boolExp,
			}
			result.Statements = append(result.Statements, statement)
		}
	}
	return result.Simplify(), p, nil
}

func ParseOperations(args []string) (Operation, error) {
	result := &CompoundStatement{}
	p := args
	for len(p) > 0 {
		switch p[0] {
		case "filter":
			op, remain, err := ParseFilters(p[1:])
			if err != nil {
				return nil, err
			}
			p = remain
			result.Statements = append(result.Statements, op)
		}
	}
	return result.Simplify(), nil
}
