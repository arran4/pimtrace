package main

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strings"
	"testing"
)

func TestFilterTokenizerScanN(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		n        int
		tokens   []any
		remainer []string
		wantErr  bool
	}{
		{
			name:     "Empty",
			args:     []string{},
			n:        1,
			tokens:   []any{},
			remainer: []string{},
			wantErr:  false,
		},
		{
			name:     "Small N does nothing",
			args:     []string{"where"},
			n:        0,
			tokens:   []any{},
			remainer: []string{"where"},
			wantErr:  false,
		},
		{
			name: "'Where' by itself",
			args: []string{"where"},
			n:    1,
			tokens: []any{
				FilterTerminator("where"),
			},
			remainer: []string{},
			wantErr:  false,
		},
		{
			name: "'Where' by itself - n in excess",
			args: []string{"where"},
			n:    10,
			tokens: []any{
				FilterTerminator("where"),
			},
			remainer: []string{},
			wantErr:  false,
		},
		{
			name: "'Where' by itself - tokens in excess",
			args: []string{"where", "where", "where", "where", "where", "where", "where", "where"},
			n:    1,
			tokens: []any{
				FilterTerminator("where"),
			},
			remainer: []string{"where", "where", "where", "where", "where", "where", "where"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := FilterTokenizerScanN(tt.args, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterTokenizerScanN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.tokens); diff != "" {
				t.Errorf("FilterTokenizerScanN() got / tt.want diff:\n %s", diff)
			}
			if !reflect.DeepEqual(got1, tt.remainer) {
				t.Errorf("FilterTokenizerScanN() got1 = %v, want %v", got1, tt.remainer)
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
				FilterTerminator("where"),
			},
			matchTokens: []any{
				FilterTerminator("where"),
			},
			want: []any{FilterTerminator("where")},
		},
		{
			name: "Terminator where and map match",
			inputTokens: []any{
				FilterTerminator("where"),
			},
			matchTokens: []any{
				FilterTerminator("map"),
			},
			want: []any{FilterTerminator("map")},
		},
		{
			name: "Terminator where and not don't match",
			inputTokens: []any{
				FilterTerminator("where"),
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
				FilterTerminator("map"),
			},
			want: nil,
		},
		{
			name: "1 tokens but no expected token types exist match",
			inputTokens: []any{
				FilterTerminator("map"),
			},
			matchTokens: []any{},
			want:        []any{},
		},
		{
			name: "1 tokens match one or the other where there is a match",
			inputTokens: []any{
				FilterTerminator("map"),
			},
			matchTokens: []any{
				[]any{
					FilterNot("not"),
					FilterTerminator("map"),
				},
			},
			want: []any{
				FilterTerminator("map"),
			},
		},
		{
			name: "1 tokens don't match one or the other where there isn't a match",
			inputTokens: []any{
				FilterTerminator("map"),
			},
			matchTokens: []any{
				[]any{
					EntryExpression("h.User-Agent"),
					FilterNot("not"),
				},
			},
			want: nil,
		},
		{
			name: "sequence of 2 match",
			inputTokens: []any{
				EntryExpression("h.User-Agent"),
				FilterNot("not"),
				FilterTerminator("map"),
			},
			matchTokens: []any{
				EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
			want: []any{
				EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
		},
		{
			name: "sequence of don't match",
			inputTokens: []any{
				EntryExpression("h.User-Agent"),
				FilterTerminator("map"),
				FilterNot("not"),
			},
			matchTokens: []any{
				EntryExpression("h.User-Agent"),
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
		statements         []Operation
		expectedExpression BooleanExpression
		remaining          []string
		wantErr            bool
	}{
		{
			name:               "Empty args go no where - since filter is already provided it's safe to die here",
			args:               []string{},
			statements:         []Operation{},
			expectedExpression: nil,
			remaining:          nil,
			wantErr:            true,
		},
		{
			name: "Basic neg expression",
			args: []string{"not", "h.user-agent", "eq", ".Kmail"},
			expectedExpression: &NotOp{
				Not: &Op{Op: EqualOp, LHS: EntryExpression("h.user-agent"), RHS: ConstantExpression("Kmail")},
			},
			statements: []Operation{},
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
			if diff := cmp.Diff(got, tt.expectedExpression, cmp.Comparer(func(o1 OpFunc, o2 OpFunc) bool {
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
		expectedOperation Operation
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
			expectedOperation: &FilterStatement{
				Expression: &NotOp{
					Not: &Op{Op: EqualOp, LHS: EntryExpression("h.user-agent"), RHS: ConstantExpression("Kmail")},
				},
			},
			remaining: []string{},
			wantErr:   false,
		},
		{
			name: "filter out into a maildir",
			args: strings.Split("filter not h.user-argent icontains .kmail into maildir", " "),
		},
		{
			name: "filter into a table",
			args: strings.Split("filter not h.user-argent icontains .kmail into table with h.user-agent h.subject year[h.date] month[h.date]", " "),
		},
		{
			name: "Filter into summary with count and a calculated sum",
			args: strings.Split("filter not h.user-argent icontains .kmail into summary count h.user-agent year[h.date] month[h.date] calculate sum[h.size]", " "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOperations(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.expectedOperation, cmp.Comparer(func(o1 OpFunc, o2 OpFunc) bool {
				sf1 := reflect.ValueOf(o1)
				sf2 := reflect.ValueOf(o2)
				return sf1.Pointer() == sf2.Pointer()
			})); diff != "" {
				t.Errorf("ParseFilter() expectedExpression %s", diff)
			}
		})
	}
}
