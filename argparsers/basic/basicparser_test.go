package basic

import (
	"github.com/google/go-cmp/cmp"
	"pimtrace/ast"
	"pimtrace/funcs"
	"reflect"
	"strings"
	"testing"
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
		expectedExpression ast.BooleanExpression
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
			expectedExpression: &ast.NotOp{
				Not: &ast.Op{Op: ast.EqualOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
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
				Expression: &ast.NotOp{
					Not: &ast.Op{Op: ast.EqualOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
				},
			},
			remaining: []string{},
			wantErr:   false,
		},
		//{
		//	name: "filter out into a mbox",
		//	args: strings.Split("filter not h.user-agent icontains .Kmail into mbox", " "),
		//	expectedOperation: &ast.CompoundStatement{
		//		Statements: []ast.Operation{
		//			&ast.FilterStatement{
		//				Expression: &ast.NotOp{
		//					Not: &ast.Op{Op: ast.IContainsOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
		//				},
		//			},
		//		},
		//	},
		//},
		{
			name: "filter into a table",
			args: strings.Split("filter not h.user-agent icontains .Kmail into table h.user-agent h.subject f.year[h.date] f.month[h.date]", " "),
			expectedOperation: &ast.CompoundStatement{
				Statements: []ast.Operation{
					&ast.FilterStatement{
						Expression: &ast.NotOp{
							Not: &ast.Op{Op: ast.IContainsOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
						},
					},
					&ast.TableTransformer{
						Columns: []*ast.ColumnExpression{
							{Name: "user-agent", Operation: ast.EntryExpression("h.user-agent")},
							{Name: "subject", Operation: ast.EntryExpression("h.subject")},
							{Name: "year-date", Operation: &ast.FunctionExpression{Function: "year", F: funcs.Year[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.date")}}},
							{Name: "month-date", Operation: &ast.FunctionExpression{Function: "month", F: funcs.Month[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.date")}}},
						},
					},
				},
			},
		},
		//{
		//	name: "filter out into a mbox sorted by date",
		//	args: strings.Split("filter not h.user-agent icontains .Kmail into mbox sort h.date", " "),
		//	expectedOperation: &ast.CompoundStatement{
		//		Statements: []ast.Operation{
		//			&ast.FilterStatement{
		//				Expression: &ast.NotOp{
		//					Not: &ast.Op{Op: ast.IContainsOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
		//				},
		//			},
		//			&maildata.MBoxOutput{},
		//			&ast.SortTransformer{
		//				Expression: []ast.ValueExpression{
		//					ast.EntryExpression("h.date"),
		//				},
		//			},
		//		},
		//	},
		//},
		{
			name: "filter into a table sorted by date",
			args: strings.Split("filter not h.user-agent icontains .Kmail into table h.user-agent h.subject f.year[h.date] f.month[h.date] sort h.date", " "),
			expectedOperation: &ast.CompoundStatement{
				Statements: []ast.Operation{
					&ast.FilterStatement{
						Expression: &ast.NotOp{
							Not: &ast.Op{Op: ast.IContainsOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
						},
					},
					&ast.TableTransformer{
						Columns: []*ast.ColumnExpression{
							{Name: "user-agent", Operation: ast.EntryExpression("h.user-agent")},
							{Name: "subject", Operation: ast.EntryExpression("h.subject")},
							{Name: "year-date", Operation: &ast.FunctionExpression{Function: "year", F: funcs.Year[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.date")}}},
							{Name: "month-date", Operation: &ast.FunctionExpression{Function: "month", F: funcs.Month[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.date")}}},
						},
					},
					&ast.SortTransformer{
						Expression: []ast.ValueExpression{
							ast.EntryExpression("h.date"),
						},
					},
				},
			},
		},
		{
			name: "Filter into summary with count and a calculated sum",
			args: strings.Split("filter not h.user-agent icontains .Kmail into summary h.user-agent h.subject f.year[h.date] f.month[h.date] calculate f.sum[h.size] f.count", " "),
			expectedOperation: &ast.CompoundStatement{
				Statements: []ast.Operation{
					&ast.FilterStatement{
						Expression: &ast.NotOp{
							Not: &ast.Op{Op: ast.IContainsOp, LHS: ast.EntryExpression("h.user-agent"), RHS: ast.ConstantExpression("Kmail")},
						},
					},
					&ast.GroupTransformer{
						Columns: []*ast.ColumnExpression{
							{Name: "user-agent", Operation: ast.EntryExpression("h.user-agent")},
							{Name: "subject", Operation: ast.EntryExpression("h.subject")},
							{Name: "year-date", Operation: &ast.FunctionExpression{Function: "year", F: funcs.Year[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.date")}}},
							{Name: "month-date", Operation: &ast.FunctionExpression{Function: "month", F: funcs.Month[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.date")}}},
						},
					},
					&ast.TableTransformer{
						Columns: []*ast.ColumnExpression{
							{Name: "user-agent", Operation: ast.EntryExpression("c.user-agent")},
							{Name: "subject", Operation: ast.EntryExpression("c.subject")},
							{Name: "year-date", Operation: ast.EntryExpression("c.year-date")},
							{Name: "month-date", Operation: ast.EntryExpression("c.month-date")},
							{Name: "sum-size", Operation: &ast.FunctionExpression{Function: "sum", F: funcs.Sum[ast.ValueExpression]{}, Args: []ast.ValueExpression{ast.EntryExpression("h.size")}}},
							{Name: "count", Operation: &ast.FunctionExpression{Function: "count", F: funcs.Count[ast.ValueExpression]{}}}, //Args: []ValueExpression{EntryExpression("t.contents")}}},
						},
					},
				},
			},
		},
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
			})); diff != "" {
				t.Errorf("ParseFilter() expectedExpression %s", diff)
			}
		})
	}
}
