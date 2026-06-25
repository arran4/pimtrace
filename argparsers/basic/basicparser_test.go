package basic

import (
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
						Expression: &evaluator.IsExpression{Field: "user-agent", Value: "Kmail"},
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
							Expression: &evaluator.IsExpression{Field: "user-agent", Value: "Kmail"},
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
									Expression: &evaluator.IContainsExpression{Field: "user-agent", Value: "Kmail"},
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
