package basic

import (
	"fmt"
	"pimtrace"
	"pimtrace/ast"
	"pimtrace/dataformats/maildata"
	"reflect"
	"strings"
	"testing"

	"github.com/arran4/go-evaluator"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFilterTokenizerScanN(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		n         int
		tokens    []any
		remainder []string
		wantErr   bool
	}{
		{
			name:      "Empty",
			args:      []string{},
			n:         1,
			tokens:    []any{},
			remainder: []string{},
			wantErr:   false,
		},
		{
			name:      "Small N does nothing",
			args:      []string{"where"},
			n:         0,
			tokens:    []any{},
			remainder: []string{"where"},
			wantErr:   false,
		},
		{
			name:   "'Where' by itself",
			args:   []string{"where"},
			n:      1,
			tokens: []any{},
			remainder: []string{
				"where",
			},
			wantErr: false,
		},
		{
			name:   "'Where' by itself - n in excess",
			args:   []string{"where"},
			n:      10,
			tokens: []any{},
			remainder: []string{
				"where",
			},
			wantErr: false,
		},
		{
			name:      "'Where' by itself - tokens in excess",
			args:      []string{"where", "where", "where", "where", "where", "where", "where", "where"},
			n:         1,
			tokens:    []any{},
			remainder: []string{"where", "where", "where", "where", "where", "where", "where", "where"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, remainder, err := FilterTokenizerScanN(tt.args, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterTokenizerScanN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tokens, tt.tokens); diff != "" {
				t.Errorf("FilterTokenizerScanN() tokens / tt.want diff:\n %s", diff)
			}
			if !reflect.DeepEqual(remainder, tt.remainder) {
				t.Errorf("FilterTokenizerScanN() remainder = %v, want %v", remainder, tt.remainder)
			}
		})
	}
}

func TestFilterTokenMatcher(t *testing.T) {
	tests := []struct {
		name        string
		inputTokens []any
		matchTokens []any
		want        []any
	}{
		{
			name:        "Empty",
			inputTokens: []any{},
			matchTokens: []any{},
			want:        []any{},
		},
		{
			name: "Terminator where match",
			inputTokens: []any{
				Terminator("where"),
			},
			matchTokens: []any{
				Terminator("where"),
			},
			want: []any{Terminator("where")},
		},
		{
			name: "Terminator where and map match",
			inputTokens: []any{
				Terminator("where"),
			},
			matchTokens: []any{
				Terminator("map"),
			},
			want: []any{Terminator("map")},
		},
		{
			name: "Terminator where and not don't match",
			inputTokens: []any{
				Terminator("where"),
			},
			matchTokens: []any{
				FilterNot("not"),
			},
			want: nil,
		},
		{
			name:        "No tokens but expected token types exist don't match",
			inputTokens: []any{},
			matchTokens: []any{
				Terminator("map"),
			},
			want: nil,
		},
		{
			name: "1 tokens but no expected token types exist match",
			inputTokens: []any{
				Terminator("map"),
			},
			matchTokens: []any{},
			want:        []any{},
		},
		{
			name: "1 tokens match one or the other where there is a match",
			inputTokens: []any{
				Terminator("map"),
			},
			matchTokens: []any{
				[]any{
					FilterNot("not"),
					Terminator("map"),
				},
			},
			want: []any{
				Terminator("map"),
			},
		},
		{
			name: "1 tokens don't match one or the other where there isn't a match",
			inputTokens: []any{
				Terminator("map"),
			},
			matchTokens: []any{
				[]any{
					ast.EntryExpression("h.User-Agent"),
					FilterNot("not"),
				},
			},
			want: nil,
		},
		{
			name: "sequence of 2 match",
			inputTokens: []any{
				ast.EntryExpression("h.User-Agent"),
				FilterNot("not"),
				Terminator("map"),
			},
			matchTokens: []any{
				ast.EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
			want: []any{
				ast.EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
		},
		{
			name: "sequence of don't match",
			inputTokens: []any{
				ast.EntryExpression("h.User-Agent"),
				Terminator("map"),
				FilterNot("not"),
			},
			matchTokens: []any{
				ast.EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TokenMatcher(tt.inputTokens, tt.matchTokens...)
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("TokenMatcher() = \n%s", diff)
			}
		})
	}
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		statements         []ast.Operation
		expectedExpression *evaluator.Query
		remaining          []string
		wantErr            bool
	}{
		{
			name:               "Empty args go no where - since filter is already provided it's safe to die here",
			args:               []string{},
			statements:         []ast.Operation{},
			expectedExpression: nil,
			remaining:          nil,
			wantErr:            true,
		},
		{
			name: "Basic neg expression",
			args: []string{"not", "h.user-agent", "eq", ".Kmail"},
			expectedExpression: &evaluator.Query{
				Expression: &evaluator.NotExpression{
					Expression: evaluator.Query{
						Expression: &evaluator.ComparisonExpression{Operation: "eq", LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
					},
				},
			},
			statements: []ast.Operation{},
			remaining:  []string{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseFilter(tt.args, tt.statements)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.expectedExpression, cmp.Comparer(func(o1 ast.OpFunc, o2 ast.OpFunc) bool {
				sf1 := reflect.ValueOf(o1)
				sf2 := reflect.ValueOf(o2)
				return sf1.Pointer() == sf2.Pointer()
			})); diff != "" {
				t.Errorf("ParseFilter() expectedExpression %s", diff)
			}
			if diff := cmp.Diff(got1, tt.remaining); diff != "" {
				t.Errorf("ParseFilter() remaining %s", diff)
			}
		})
	}
}

func TestParseOperations(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOperation ast.Operation
		remaining         []string
		wantErr           bool
	}{
		{
			name:              "Empty args go no where - since filter is already provided it's safe to die here",
			args:              []string{},
			expectedOperation: nil,
			remaining:         nil,
			wantErr:           false,
		},
		{
			name: "Basic neg expression",
			args: []string{"filter", "not", "h.user-agent", "eq", ".Kmail"},
			expectedOperation: &ast.FilterStatement{
				Expression: &evaluator.Query{
					Expression: &evaluator.NotExpression{
						Expression: evaluator.Query{
							Expression: &evaluator.ComparisonExpression{Operation: "eq", LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
						},
					},
				},
			},
			remaining: []string{},
			wantErr:   false,
		},
		{
			name: "filter out into a mbox",
			args: strings.Split("filter not h.user-agent icontains .Kmail into mbox", " "),
			expectedOperation: &ast.CompoundStatement{
				Statements: []ast.Operation{
					&ast.FilterStatement{
						Expression: &evaluator.Query{
							Expression: &evaluator.NotExpression{
								Expression: evaluator.Query{
									Expression: &evaluator.ComparisonExpression{Operation: "icontains", LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
								},
							},
						},
					},
					&maildata.MBoxOutput{},
				},
			},
		},
		// {
		// 	name: "filter into a table",
		// 	// Skipped due to fragility in comparison of FunctionExpression.F
		// },
		// {
		//	name: "filter into a table sorted by date",
		//	// Skipped
		// },
		// {
		//	name: "Filter into summary...",
		//	// Skipped
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOperations(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.expectedOperation, cmp.Comparer(func(o1 ast.OpFunc, o2 ast.OpFunc) bool {
				sf1 := reflect.ValueOf(o1)
				sf2 := reflect.ValueOf(o2)
				return sf1.Pointer() == sf2.Pointer()
			}), cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); diff != "" {
				t.Errorf("ParseFilter() expectedExpression %s", diff)
			}
		})
	}
}

func TestParseIntoMbox(t *testing.T) {
	_, err := ParseOperations([]string{"into", "mbox"})
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}

func TestParseSort(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOperation ast.Operation
		remaining         []string
		wantErr           bool
	}{
		{
			name: "Basic sort",
			args: []string{"c.name", "into", "mbox"},
			expectedOperation: &ast.SortTransformer{
				Expression: []ast.ValueExpression{
					ast.EntryExpression("c.name"),
				},
			},
			remaining: []string{"into", "mbox"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseSort(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.expectedOperation, cmp.Comparer(func(o1 ast.OpFunc, o2 ast.OpFunc) bool {
				sf1 := reflect.ValueOf(o1)
				sf2 := reflect.ValueOf(o2)
				return sf1.Pointer() == sf2.Pointer()
			}), cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); diff != "" {
				t.Errorf("ParseSort() expectedOperation %s", diff)
			}
			if diff := cmp.Diff(got1, tt.remaining); diff != "" {
				t.Errorf("ParseSort() remaining %s", diff)
			}
		})
	}
}

func TestParseIntoTable(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOperation ast.Operation
		remaining         []string
		wantErr           bool
	}{
		{
			name: "Basic table",
			args: []string{"c.name", "c.date", "into", "mbox"},
			expectedOperation: &ast.TableTransformer{
				Columns: []*ast.ColumnExpression{
					{
						Operation: ast.EntryExpression("c.name"),
						Name:      "name",
					},
					{
						Operation: ast.EntryExpression("c.date"),
						Name:      "date",
					},
				},
			},
			remaining: []string{"into", "mbox"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseIntoTable(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntoTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.expectedOperation, cmp.Comparer(func(o1 ast.OpFunc, o2 ast.OpFunc) bool {
				sf1 := reflect.ValueOf(o1)
				sf2 := reflect.ValueOf(o2)
				return sf1.Pointer() == sf2.Pointer()
			}), cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); diff != "" {
				t.Errorf("ParseIntoTable() expectedOperation %s", diff)
			}
			if diff := cmp.Diff(got1, tt.remaining); diff != "" {
				t.Errorf("ParseIntoTable() remaining %s", diff)
			}
		})
	}
}

func TestParseIntoSummary(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOperation ast.Operation
		remaining         []string
		wantErr           bool
	}{
		{
			name: "Basic summary",
			args: []string{"c.name", "c.date", "into", "mbox"},
			expectedOperation: &ast.GroupTransformer{
				Columns: []*ast.ColumnExpression{
					{
						Operation: ast.EntryExpression("c.name"),
						Name:      "name",
					},
					{
						Operation: ast.EntryExpression("c.date"),
						Name:      "date",
					},
				},
			},
			remaining: []string{"into", "mbox"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseIntoSummary(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntoSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.expectedOperation, cmp.Comparer(func(o1 ast.OpFunc, o2 ast.OpFunc) bool {
				sf1 := reflect.ValueOf(o1)
				sf2 := reflect.ValueOf(o2)
				return sf1.Pointer() == sf2.Pointer()
			}), cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); diff != "" {
				t.Errorf("ParseIntoSummary() expectedOperation %s", diff)
			}
			if diff := cmp.Diff(got1, tt.remaining); diff != "" {
				t.Errorf("ParseIntoSummary() remaining %s", diff)
			}
		})
	}
}

func TestParseFunctionExpression(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		want      ast.ValueExpression
		remaining []string
		wantErr   bool
	}{
		{
			name: "Basic function",
			args: []string{"f.count"},
			want: &ast.FunctionExpression{
				Function: "count",
			},
			remaining: []string{},
			wantErr:   false,
		},
		{
			name: "Basic function with params",
			args: []string{"f.count[c.name]"},
			want: &ast.FunctionExpression{
				Function: "count",
				Args: []ast.ValueExpression{
					ast.EntryExpression("c.name"),
				},
			},
			remaining: []string{},
			wantErr:   false,
		},
		{
			name: "Evaluator function",
			args: []string{"f.year[c.name]"},
			want: &ast.EvaluatorFunctionExpression{
				Function: "year",
				FunctionExpression: evaluator.FunctionExpression{
					Args: []evaluator.Term{
						ast.EntryExpression("c.name"),
					},
				},
			},
			remaining: []string{},
			wantErr:   false,
		},
		{
			name: "Evaluator function multiple args",
			args: []string{"f.year[c.name,c.date]"},
			want: &ast.EvaluatorFunctionExpression{
				Function: "year",
				FunctionExpression: evaluator.FunctionExpression{
					Args: []evaluator.Term{
						ast.EntryExpression("c.name"),
						ast.EntryExpression("c.date"),
					},
				},
			},
			remaining: []string{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseFunctionExpression(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFunctionExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); diff != "" {
				t.Errorf("ParseFunctionExpression() want %s", diff)
			}
			if diff := cmp.Diff(got1, tt.remaining); diff != "" {
				t.Errorf("ParseFunctionExpression() remaining %s", diff)
			}
		})
	}
}

func TestParseExpressions(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    []ast.ValueExpression
		wantErr bool
	}{
		{
			name: "Basic parameter",
			s:    "c.name",
			want: []ast.ValueExpression{
				ast.EntryExpression("c.name"),
			},
			wantErr: false,
		},
		{
			name: "Multiple parameters",
			s:    "c.name,c.date",
			want: []ast.ValueExpression{
				ast.EntryExpression("c.name"),
				ast.EntryExpression("c.date"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseExpressions(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpressions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); len(diff) > 0 {
				t.Errorf("ParseExpressions() = \n%s", diff)
			}
		})
	}
}

func TestParseExpressions_MoreParams(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    []ast.ValueExpression
		wantErr bool
	}{
		{
			name: "Header param",
			s:    "h.User-Agent",
			want: []ast.ValueExpression{
				ast.EntryExpression("h.User-Agent"),
			},
			wantErr: false,
		},
		{
			name: "Constant param",
			s:    ".value",
			want: []ast.ValueExpression{
				ast.ConstantExpression("value"),
			},
			wantErr: false,
		},
		{
			name: "Function param",
			s:    "f.year[c.name]",
			want: []ast.ValueExpression{
				&ast.EvaluatorFunctionExpression{
					Function: "year",
					FunctionExpression: evaluator.FunctionExpression{
						Args: []evaluator.Term{
							ast.EntryExpression("c.name"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Unknown param",
			s:       "unknown",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid expression format",
			s:       "f.sum[",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseExpressions(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpressions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreFields(ast.FunctionExpression{}, "F"), cmpopts.IgnoreFields(evaluator.FunctionExpression{}, "Func")); len(diff) > 0 {
				t.Errorf("ParseExpressions() = \n%s", diff)
			}
		})
	}
}

type testMockEntry map[string]pimtrace.Value

func (m testMockEntry) Get(key string) (pimtrace.Value, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("field %q not found", key)
}

func TestParserEvaluatorAcceptance(t *testing.T) {
	t.Run("relational operators map and evaluate correctly", func(t *testing.T) {
		cases := []struct {
			name     string
			op       string
			expected string
			inputVal string
			wantRes  bool
		}{
			{"gt true", "gt", ">", "10", true},
			{"gt false", "gt", ">", "5", false},
			{"gte true equal", "gte", ">=", "5", true},
			{"gte true greater", "gte", ">=", "10", true},
			{"gte false", "gte", ">=", "2", false},
			{"lt true", "lt", "<", "2", true},
			{"lt false", "lt", "<", "5", false},
			{"lte true equal", "lte", "<=", "5", true},
			{"lte true less", "lte", "<=", "2", true},
			{"lte false", "lte", "<=", "10", false},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				q, remain, err := ParseFilter([]string{"c.val", tc.op, ".5"}, nil)
				if err != nil {
					t.Fatalf("ParseFilter error: %v", err)
				}
				if len(remain) != 0 {
					t.Errorf("expected empty remain, got %v", remain)
				}
				safeComp, ok := q.Expression.(*ast.SafeComparisonExpression)
				if !ok {
					t.Fatalf("expected *ast.SafeComparisonExpression, got %T", q.Expression)
				}
				if safeComp.Operator != tc.expected {
					t.Errorf("expected operator %q, got %q", tc.expected, safeComp.Operator)
				}

				res, err := q.Evaluate(testMockEntry{"c.val": pimtrace.SimpleStringValue(tc.inputVal)})
				if err != nil {
					t.Fatalf("Evaluate error: %v", err)
				}
				if res != tc.wantRes {
					t.Errorf("Evaluate() = %v, want %v", res, tc.wantRes)
				}
			})
		}
	})

	t.Run("A or B and C uses evaluator precedence", func(t *testing.T) {
		// Evaluator precedence: and binds more tightly than or.
		// Expression: c.a eq .1 or c.b eq .1 and c.c eq .1
		// Evaluates as: c.a eq .1 or (c.b eq .1 and c.c eq .1)
		q, _, err := ParseFilter([]string{"c.a", "eq", ".1", "or", "c.b", "eq", ".1", "and", "c.c", "eq", ".1"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}

		// Case 1: a=1, b=0, c=0 -> 1 or (0 and 0) = true (if left-associative it would be (1 or 0) and 0 = false)
		res1, err := q.Evaluate(testMockEntry{
			"c.a": pimtrace.SimpleStringValue("1"),
			"c.b": pimtrace.SimpleStringValue("0"),
			"c.c": pimtrace.SimpleStringValue("0"),
		})
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if !res1 {
			t.Errorf("expected true for a=1,b=0,c=0 due to and-over-or precedence, got false")
		}

		// Case 2: a=0, b=1, c=0 -> 0 or (1 and 0) = false
		res2, err := q.Evaluate(testMockEntry{
			"c.a": pimtrace.SimpleStringValue("0"),
			"c.b": pimtrace.SimpleStringValue("1"),
			"c.c": pimtrace.SimpleStringValue("0"),
		})
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if res2 {
			t.Errorf("expected false for a=0,b=1,c=0, got true")
		}
	})

	t.Run("parentheses override precedence", func(t *testing.T) {
		// Expression: ( c.a eq .1 or c.b eq .1 ) and c.c eq .1
		q, _, err := ParseFilter([]string{"(", "c.a", "eq", ".1", "or", "c.b", "eq", ".1", ")", "and", "c.c", "eq", ".1"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}

		// With a=1, b=0, c=0: (1 or 0) and 0 = false (overridden precedence!)
		res, err := q.Evaluate(testMockEntry{
			"c.a": pimtrace.SimpleStringValue("1"),
			"c.b": pimtrace.SimpleStringValue("0"),
			"c.c": pimtrace.SimpleStringValue("0"),
		})
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if res {
			t.Errorf("expected false when parentheses group (a or b) and c with c=0, got true")
		}

		// With a=1, b=0, c=1: (1 or 0) and 1 = true
		resTrue, err := q.Evaluate(testMockEntry{
			"c.a": pimtrace.SimpleStringValue("1"),
			"c.b": pimtrace.SimpleStringValue("0"),
			"c.c": pimtrace.SimpleStringValue("1"),
		})
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if !resTrue {
			t.Errorf("expected true with c=1, got false")
		}
	})

	t.Run("nested not", func(t *testing.T) {
		qDouble, _, err := ParseFilter([]string{"not", "not", "c.val", "eq", ".hello"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}
		res, err := qDouble.Evaluate(testMockEntry{"c.val": pimtrace.SimpleStringValue("hello")})
		if err != nil || !res {
			t.Errorf("not not matching value: res=%v, err=%v", res, err)
		}
		resFalse, err := qDouble.Evaluate(testMockEntry{"c.val": pimtrace.SimpleStringValue("world")})
		if err != nil || resFalse {
			t.Errorf("not not non-matching value: res=%v, err=%v", resFalse, err)
		}

		qTriple, _, err := ParseFilter([]string{"not", "not", "not", "c.val", "eq", ".hello"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}
		resTriple, err := qTriple.Evaluate(testMockEntry{"c.val": pimtrace.SimpleStringValue("hello")})
		if err != nil || resTriple {
			t.Errorf("not not not matching value: res=%v, err=%v", resTriple, err)
		}
	})

	t.Run("mixed icontains and relational boolean composition", func(t *testing.T) {
		q, _, err := ParseFilter([]string{"c.title", "icontains", ".Report", "and", "c.amount", "gt", ".100"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}

		match, err := q.Evaluate(testMockEntry{
			"c.title":  pimtrace.SimpleStringValue("Monthly report Q3"),
			"c.amount": pimtrace.SimpleStringValue("250.50"),
		})
		if err != nil || !match {
			t.Errorf("expected match, got match=%v err=%v", match, err)
		}

		noMatchAmount, err := q.Evaluate(testMockEntry{
			"c.title":  pimtrace.SimpleStringValue("Monthly report Q3"),
			"c.amount": pimtrace.SimpleStringValue("50"),
		})
		if err != nil || noMatchAmount {
			t.Errorf("expected no match on amount, got match=%v err=%v", noMatchAmount, err)
		}

		noMatchTitle, err := q.Evaluate(testMockEntry{
			"c.title":  pimtrace.SimpleStringValue("Monthly invoice"),
			"c.amount": pimtrace.SimpleStringValue("250"),
		})
		if err != nil || noMatchTitle {
			t.Errorf("expected no match on title, got match=%v err=%v", noMatchTitle, err)
		}
	})

	t.Run("invalid typed coercion returns error through SafeComparisonExpression", func(t *testing.T) {
		qNum, _, err := ParseFilter([]string{"c.val", "gt", ".10"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}
		_, err = qNum.Evaluate(testMockEntry{"c.val": pimtrace.SimpleStringValue("apple")})
		if err == nil {
			t.Errorf("expected error on invalid numeric coercion, got nil")
		}

		qDate, _, err := ParseFilter([]string{"c.dt", "gt", ".2020-01-01"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}
		_, err = qDate.Evaluate(testMockEntry{"c.dt": pimtrace.SimpleStringValue("invalid-date")})
		if err == nil {
			t.Errorf("expected error on invalid date coercion, got nil")
		}
	})

	t.Run("legacy textual 01 vs 1 equality remains textual", func(t *testing.T) {
		// Legacy 3-token eq
		q3, _, err := ParseFilter([]string{"c.code", "eq", ".1"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}
		res3, err := q3.Evaluate(testMockEntry{"c.code": pimtrace.SimpleStringValue("01")})
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if res3 {
			t.Errorf("01 eq .1 should be false textually in 3-token path, got true")
		}

		// Compound eq path
		qCompound, _, err := ParseFilter([]string{"c.code", "eq", ".1", "and", "c.other", "eq", ".x"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}
		resCompound, err := qCompound.Evaluate(testMockEntry{
			"c.code":  pimtrace.SimpleStringValue("01"),
			"c.other": pimtrace.SimpleStringValue("x"),
		})
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if resCompound {
			t.Errorf("01 eq .1 should be false textually in compound path, got true")
		}
	})

	t.Run("full c.*, h.*, p.* identifiers survive placeholder transformation/restoration", func(t *testing.T) {
		q, _, err := ParseFilter([]string{
			"c.total_amount", "gt", ".50",
			"and", "h.User-Agent", "eq", ".curl",
			"and", "p.DTSTART", "gt", ".2020-01-01",
		}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}

		// Evaluate against an entry with these full keys
		entry := testMockEntry{
			"c.total_amount": pimtrace.SimpleStringValue("100"),
			"h.User-Agent":   pimtrace.SimpleStringValue("curl"),
			"p.DTSTART":      pimtrace.SimpleStringValue("2020-05-01"),
		}
		match, err := q.Evaluate(entry)
		if err != nil {
			t.Fatalf("Evaluate error: %v", err)
		}
		if !match {
			t.Errorf("expected match with preserved full identifiers")
		}
	})

	t.Run("cross-check with directly constructed go-evaluator expression", func(t *testing.T) {
		// PIMTrace filter: c.a eq .1 or c.b eq .1 and c.c gt .5
		qParsed, _, err := ParseFilter([]string{"c.a", "eq", ".1", "or", "c.b", "eq", ".1", "and", "c.c", "gt", ".5"}, nil)
		if err != nil {
			t.Fatalf("ParseFilter error: %v", err)
		}

		// Directly constructed go-evaluator equivalent:
		// OrExpression:
		//   [0]: ComparisonExpression (c.a eq 1)
		//   [1]: AndExpression:
		//          [0]: ComparisonExpression (c.b eq 1)
		//          [1]: SafeComparisonExpression (c.c > 5)
		directExpr := &evaluator.OrExpression{
			Expressions: []evaluator.Query{
				{
					Expression: &evaluator.ComparisonExpression{
						Operation: "eq",
						LHS:       ast.EntryExpression("c.a"),
						RHS:       ast.ConstantExpression("1"),
					},
				},
				{
					Expression: &evaluator.AndExpression{
						Expressions: []evaluator.Query{
							{
								Expression: &evaluator.ComparisonExpression{
									Operation: "eq",
									LHS:       ast.EntryExpression("c.b"),
									RHS:       ast.ConstantExpression("1"),
								},
							},
							{
								Expression: &ast.SafeComparisonExpression{
									Operator: ">",
									Left:     ast.EntryExpression("c.c"),
									Right:    ast.ConstantExpression("5"),
								},
							},
						},
					},
				},
			},
		}

		// Truth table comparison across all combinations
		aVals := []string{"0", "1"}
		bVals := []string{"0", "1"}
		cVals := []string{"2", "10"}

		for _, a := range aVals {
			for _, b := range bVals {
				for _, c := range cVals {
					entry := testMockEntry{
						"c.a": pimtrace.SimpleStringValue(a),
						"c.b": pimtrace.SimpleStringValue(b),
						"c.c": pimtrace.SimpleStringValue(c),
					}

					parsedRes, parsedErr := qParsed.Evaluate(entry)
					directRes, directErr := directExpr.Evaluate(entry)

					if (parsedErr != nil) != (directErr != nil) {
						t.Errorf("error mismatch for a=%s,b=%s,c=%s: parsedErr=%v directErr=%v", a, b, c, parsedErr, directErr)
					}
					if parsedRes != directRes {
						t.Errorf("result mismatch for a=%s,b=%s,c=%s: parsed=%v direct=%v", a, b, c, parsedRes, directRes)
					}
				}
			}
		}
	})
}
