package main

import (
	"fmt"
	"reflect"
	"strings"
)

var ErrParserNothingFound = fmt.Errorf("no token found")
var ErrUnknownExpression = fmt.Errorf("unknown expression")
var ErrParserFault = fmt.Errorf("parser fault")
var ErrUnknownIntoStatement = fmt.Errorf("unknown into")

type FilterEquals string
type FilterContains string
type FilterIContains string
type FilterNot string
type FilterTerminator string

var _ ValueExpression = ConstantExpression("")
var _ ValueExpression = EntryExpression("")

func FilterIdentify(s string) (any, error) {
	ss := strings.SplitN(s, ".", 2)
	switch ss[0] {
	case "into", "filter", "where", "sort":
		return FilterTerminator(s), nil
	case "not":
		return FilterNot(s), nil
	case "eq":
		return FilterEquals(s), nil
	case "contains":
		return FilterContains(s), nil
	case "icontains":
		return FilterIContains(s), nil
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

func ParseFilter(args []string, statements []Operation) (BooleanExpression, []string, error) {
	tks, remain, err := FilterTokenizerScanN(args, 3)
	if err != nil {
		return nil, nil, err
	}
	if TokenMatcher(tks, FilterNot("")) != nil {
		var op BooleanExpression
		op, remain, err = ParseFilter(args[1:], []Operation{})
		if err != nil {
			return nil, nil, err
		}
		return &NotOp{
			Not: op,
		}, remain, nil
	}
	if matches := TokenMatcher(tks,
		[]any{EntryExpression(""), ConstantExpression("")},
		[]any{FilterEquals(""), FilterContains(""), FilterIContains("")},
		[]any{EntryExpression(""), ConstantExpression("")},
	); len(matches) > 0 {
		var op OpFunc
		switch /*opMatch :=*/ matches[1].(type) {
		case FilterEquals:
			op = EqualOp
		case FilterContains:
			op = ContainsOp
		case FilterIContains:
			op = IContainsOp
		}
		return &Op{
			Op:  op,
			LHS: tks[0].(ValueExpression),
			RHS: tks[2].(ValueExpression),
		}, remain, nil
	}

	return nil, nil, fmt.Errorf("at %v: %w", tks, ErrParserNothingFound)
}

func TokenMatcher(inputTokens []any, matchTokens ...any) []any {
	var result []any = nil
	for i := 0; i < len(matchTokens); i++ {
		var m any = nil
		if i >= len(inputTokens) && i < len(matchTokens) {
			return nil
		}
		switch reflect.TypeOf(matchTokens[i]).Kind() {
		case reflect.Slice:
			for _, stt := range (matchTokens[i]).([]any) {
				if reflect.TypeOf(inputTokens[i]) == reflect.TypeOf(stt) {
					stt := stt
					m = stt
					break
				}
			}
			if m == nil {
				return result
			}
		case reflect.String:
			if reflect.TypeOf(inputTokens[i]) != reflect.TypeOf(matchTokens[i]) {
				return nil
			}
			m = matchTokens[i]
		}
		result = append(result, m)
	}
	if result == nil {
		result = make([]any, 0)
	}
	return result
}

func ParseFilters(args []string) (Operation, []string, error) {
	result := &CompoundStatement{}
	p := args
	for len(p) > 0 {
		switch p[0] {
		case "into":
			return result.Simplify(), p, nil
		case "filter", "where":
			p = p[1:]
			fallthrough
		default:
			boolExp, remain, err := ParseFilter(p[:], result.Statements)
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

func ParseInto(args []string) (Operation, []string, error) {
	p := args
	if len(p) > 0 {
		switch p[0] {
		case "into":
			p = p[1:]
		}
	}
	if len(p) > 0 {
		switch p[0] {
		case "mbox":
			return &MBoxOutput{}, p[1:], nil
		case "summary":
			return ParseIntoSummary(p[:], result.Statements)
		case "table":
			return ParseIntoTable(p[:], result.Statements)
		}
	}
	return nil, nil, ErrUnknownIntoStatement
}

func ParseOperations(args []string) (Operation, error) {
	result := &CompoundStatement{}
	p := args
	for len(p) > 0 {
		l := len(p)
		switch p[0] {
		case "filter":
			op, remain, err := ParseFilters(p[1:])
			if err != nil {
				return nil, err
			}
			p = remain
			if op != nil {
				result.Statements = append(result.Statements, op)
			}
		case "into":
			op, remain, err := ParseInto(p[1:])
			if err != nil {
				return nil, err
			}
			p = remain
			if op != nil {
				result.Statements = append(result.Statements, op)
			}
		case "sort":
			op, remain, err := ParseSort(p[1:])
			if err != nil {
				return nil, err
			}
			p = remain
			if op != nil {
				result.Statements = append(result.Statements, op)
			}
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownExpression, p[0])
		}
		if len(p) == l {
			return nil, ErrParserFault
		}
	}
	return result.Simplify(), nil
}
