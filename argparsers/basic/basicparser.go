package basic

import (
	"fmt"
	"pimtrace/ast"
	"reflect"
	"regexp"
	"strings"
)

var (
	ErrParserNothingFound        = fmt.Errorf("no token found")
	ErrParserUnknownToken        = fmt.Errorf("unknown token")
	ErrUnknownExpression         = fmt.Errorf("unknown expression")
	ErrParserFault               = fmt.Errorf("parser fault")
	ErrUnknownIntoStatement      = fmt.Errorf("unknown into")
	ErrIntoNotImplemented        = fmt.Errorf("into not implemented")
	ErrInvalidFunctionExpression = fmt.Errorf("invalid function expression")
)

type FilterEquals string
type FilterContains string
type FilterIContains string
type FilterNot string
type Terminator string

var _ ast.ValueExpression = ast.ConstantExpression("")
var _ ast.ValueExpression = ast.EntryExpression("")

func FilterIdentify(s string) (any, error) {
	ss := strings.SplitN(s, ".", 2)
	switch ss[0] {
	case "into", "filter", "where", "sort":
		return Terminator(s), nil
	case "not":
		return FilterNot(s), nil
	case "eq":
		return FilterEquals(s), nil
	case "contains":
		return FilterContains(s), nil
	case "icontains":
		return FilterIContains(s), nil
	case "h", "header":
		return ast.EntryExpression(s), nil
	case "c", "column":
		return ast.EntryExpression(s), nil
	case "":
		if strings.HasPrefix(s, ".") {
			return ast.ConstantExpression(ss[1]), nil
		}
	}
	return nil, fmt.Errorf("filter tokenizer: %w: %s", ErrParserUnknownToken, ss[0])
}

func IntoIdentify(args []string) (any, []string, error) {
	if len(args) == 0 {
		return nil, args, nil
	}
	ss := strings.SplitN(args[0], ".", 2)
	switch ss[0] {
	case "into", "filter", "where", "sort", "calculate":
		return Terminator(args[0]), args[0:], nil
	case "h", "header":
		return ast.EntryExpression(args[0]), args[1:], nil
	case "c", "column":
		return ast.EntryExpression(args[0]), args[1:], nil
	case "f", "func":
		return ParseFunctionExpression(args)
	case "":
		if strings.HasPrefix(args[0], ".") {
			return ast.ConstantExpression(ss[1]), args[1:], nil
		}
	}
	return nil, nil, fmt.Errorf("into tokenizer: %w: %s", ErrParserUnknownToken, ss[0])
}

func FunctionParameterExpressionIdentify(args []string) (any, []string, error) {
	if len(args) == 0 {
		return nil, args, nil
	}
	ss := strings.SplitN(args[0], ".", 2)
	switch ss[0] {
	case "h", "header":
		return ast.EntryExpression(args[0]), args[1:], nil
	case "c", "column":
		return ast.EntryExpression(args[0]), args[1:], nil
	case "f", "func":
		return ParseFunctionExpression(args)
	case "":
		if strings.HasPrefix(args[0], ".") {
			return ast.ConstantExpression(ss[1]), args[1:], nil
		}
	}
	return nil, nil, fmt.Errorf("function param tokenizer: %w: %s", ErrParserUnknownToken, ss[0])
}

var fere = regexp.MustCompile(`^(f|func)\.([^[]+)(\[([^]]+)\])?$`)

func ParseFunctionExpression(args []string) (ast.ValueExpression, []string, error) {
	m := fere.FindStringSubmatch(args[0])
	if len(m) == 5 {
		var params []ast.ValueExpression
		if len(m[3]) > 0 {
			var err error
			params, err = ParseExpressions(m[4])
			if err != nil {
				return nil, nil, fmt.Errorf("parameter parse error: %w", err)
			}
		}
		return &ast.FunctionExpression{
			Function: m[2],
			Args:     params,
		}, args[1:], nil
	}
	return nil, nil, fmt.Errorf("%w: %s", ErrInvalidFunctionExpression, args[0])
}

func ParseExpressions(s string) ([]ast.ValueExpression, error) {
	css := strings.SplitN(s, ",", 2)
	var results []ast.ValueExpression
	for len(css) > 0 {
		var err error
		var v any
		v, css, err = FunctionParameterExpressionIdentify(css)
		if err != nil {
			return nil, err
		}
		switch v := v.(type) {
		case ast.ValueExpression:
			results = append(results, v)
		default:
			return nil, fmt.Errorf("%w: %s", ErrParserUnknownToken, css[0])
		}
	}
	return results, nil
}

func FilterTokenizerScanN(args []string, n int) ([]any, []string, error) {
	i := 0
	r := []any{}
done:
	for ; i < n && i < len(args); i++ {
		t, err := FilterIdentify(args[i])
		if err != nil {
			return nil, nil, err
		}
		switch t.(type) {
		case Terminator:
			break done
		}
		r = append(r, t)
	}
	return r, args[i:], nil
}

func IntoTokenizerScan(args []string) ([]any, []string, error) {
	var r []any
done:
	for len(args) > 0 {
		var err error
		var t any
		t, args, err = IntoIdentify(args)
		if err != nil {
			return nil, nil, err
		}
		switch t.(type) {
		case Terminator:
			break done
		}
		r = append(r, t)
	}
	return r, args, nil
}

func ParseFilter(args []string, statements []ast.Operation) (ast.BooleanExpression, []string, error) {
	tks, remain, err := FilterTokenizerScanN(args, 3)
	if err != nil {
		return nil, nil, err
	}
	if TokenMatcher(tks, FilterNot("")) != nil {
		var op ast.BooleanExpression
		op, remain, err = ParseFilter(args[1:], []ast.Operation{})
		if err != nil {
			return nil, nil, err
		}
		return &ast.NotOp{
			Not: op,
		}, remain, nil
	}
	if matches := TokenMatcher(tks,
		[]any{ast.EntryExpression(""), ast.ConstantExpression("")},
		[]any{FilterEquals(""), FilterContains(""), FilterIContains("")},
		[]any{ast.EntryExpression(""), ast.ConstantExpression("")},
	); len(matches) > 1 {
		var op ast.OpFunc
		switch /*opMatch :=*/ matches[1].(type) {
		case FilterEquals:
			op = ast.EqualOp
		case FilterContains:
			op = ast.ContainsOp
		case FilterIContains:
			op = ast.IContainsOp
		}
		return &ast.Op{
			Op:  op,
			LHS: tks[0].(ast.ValueExpression),
			RHS: tks[2].(ast.ValueExpression),
		}, remain, nil
	}
	return nil, nil, fmt.Errorf("at %v: %w", tks, ErrParserNothingFound)
}

func ParseIntoSummary(args []string) (ast.Operation, []string, error) {
	results, remain, err := ParseIntoTable(args)
	if err != nil {
		return nil, nil, fmt.Errorf("summary table: %w", err)
	}
	table := &ast.GroupTransformer{
		Columns: results.(*ast.TableTransformer).Columns,
	}
	if len(remain) > 0 {
		switch remain[0] {
		case "calculate":
			remain = remain[1:]
		}
		c := &ast.CompoundStatement{
			Statements: []ast.Operation{
				table,
			},
		}
		var tks []any
		tks, remain, err = IntoTokenizerScan(remain)
		if err != nil {
			return nil, nil, fmt.Errorf("summary table: %w", err)
		}
		t := &ast.TableTransformer{}
		for _, origC := range table.Columns {
			t.Columns = append(t.Columns, &ast.ColumnExpression{
				Name:      origC.Name,
				Operation: ast.EntryExpression("c." + origC.Name),
			})
		}
	done:
		for _, tkn := range tks {
			switch tkn := tkn.(type) {
			case ast.ValueExpression:
				t.Columns = append(t.Columns, &ast.ColumnExpression{
					Name:      tkn.ColumnName(),
					Operation: tkn,
				})
			case Terminator:
				break done
			default:
				return nil, nil, fmt.Errorf("at %v: %w: unexpected token type %s", tks, ErrParserFault, reflect.TypeOf(tkn))
			}
		}
		if len(t.Columns) > len(table.Columns) {
			c.Statements = append(c.Statements, t)
		}
		results = c.Simplify()
	}
	return results, remain, nil
}

func ParseIntoTable(args []string) (ast.Operation, []string, error) {
	tks, remain, err := IntoTokenizerScan(args)
	if err != nil {
		return nil, nil, fmt.Errorf("table: %w", err)
	}
	if len(tks) > 0 {
		var expressions []*ast.ColumnExpression
	done:
		for _, tkn := range tks {
			switch tkn := tkn.(type) {
			case ast.ValueExpression:
				expressions = append(expressions, &ast.ColumnExpression{
					Name:      tkn.ColumnName(),
					Operation: tkn,
				})
			case Terminator:
				break done
			default:
				return nil, nil, fmt.Errorf("at %v: %w: unexpected token type %s", tks, ErrParserFault, reflect.TypeOf(tkn))
			}
		}
		result := &ast.TableTransformer{
			Columns: expressions,
		}
		return result, remain, nil
	}
	return nil, nil, fmt.Errorf("at %v: %w", tks, ErrParserNothingFound)
}

func ParseSort(args []string) (ast.Operation, []string, error) {
	tks, remain, err := IntoTokenizerScan(args)
	if err != nil {
		return nil, nil, err
	}
	if len(tks) > 0 {
		var expressions []ast.ValueExpression
	done:
		for _, tkn := range tks {
			switch tkn := tkn.(type) {
			case ast.ValueExpression:
				expressions = append(expressions, tkn)
			case Terminator:
				break done
			default:
				return nil, nil, fmt.Errorf("at %v: %w: unexpected token type %s", tks, ErrParserFault, reflect.TypeOf(tkn))
			}
		}
		result := &ast.SortTransformer{
			Expression: expressions,
		}
		return result, remain, nil
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

func ParseFilters(args []string) (ast.Operation, []string, error) {
	result := &ast.CompoundStatement{}
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
			statement := &ast.FilterStatement{
				Expression: boolExp,
			}
			result.Statements = append(result.Statements, statement)
		}
	}
	return result.Simplify(), p, nil
}

func ParseInto(args []string) (ast.Operation, []string, error) {
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
			//return &maildata.MBoxOutput{}, p[1:], nil
			return nil, nil, fmt.Errorf("%w: mbox", ErrIntoNotImplemented)
		case "summary":
			return ParseIntoSummary(p[1:])
		case "table":
			return ParseIntoTable(p[1:])
		}
	}
	return nil, nil, ErrUnknownIntoStatement
}

func ParseOperations(args []string) (ast.Operation, error) {
	result := &ast.CompoundStatement{}
	p := args
	for len(p) > 0 {
		l := len(p)
		switch p[0] {
		case "filter":
			op, remain, err := ParseFilters(p[1:])
			if err != nil {
				return nil, fmt.Errorf("parse filters: %w", err)
			}
			p = remain
			if op != nil {
				result.Statements = append(result.Statements, op)
			}
		case "into":
			op, remain, err := ParseInto(p[1:])
			if err != nil {
				return nil, fmt.Errorf("parse into: %w", err)
			}
			p = remain
			if op != nil {
				result.Statements = append(result.Statements, op)
			}
		case "sort":
			op, remain, err := ParseSort(p[1:])
			if err != nil {
				return nil, fmt.Errorf("parse sort: %w", err)
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
