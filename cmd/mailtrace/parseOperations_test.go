package main

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
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
		name       string
		tks        []any
		tokenTypes []any
		want       bool
	}{
		{
			name:       "Empty",
			tks:        []any{},
			tokenTypes: []any{},
			want:       true,
		},
		{
			name: "Terminator where match",
			tks: []any{
				FilterTerminator("where"),
			},
			tokenTypes: []any{
				FilterTerminator("where"),
			},
			want: true,
		},
		{
			name: "Terminator where and map match",
			tks: []any{
				FilterTerminator("where"),
			},
			tokenTypes: []any{
				FilterTerminator("map"),
			},
			want: true,
		},
		{
			name: "Terminator where and not don't match",
			tks: []any{
				FilterTerminator("where"),
			},
			tokenTypes: []any{
				FilterNot("not"),
			},
			want: false,
		},
		{
			name: "No tokens but expected token types exist don't match",
			tks:  []any{},
			tokenTypes: []any{
				FilterTerminator("map"),
			},
			want: false,
		},
		{
			name: "1 tokens but no expected token types exist match",
			tks: []any{
				FilterTerminator("map"),
			},
			tokenTypes: []any{},
			want:       true,
		},
		{
			name: "1 tokens match one or the other where there is a match",
			tks: []any{
				FilterTerminator("map"),
			},
			tokenTypes: []any{
				[]any{
					FilterNot("not"),
					FilterTerminator("map"),
				},
			},
			want: true,
		},
		{
			name: "1 tokens don't match one or the other where there isn't a match",
			tks: []any{
				FilterTerminator("map"),
			},
			tokenTypes: []any{
				[]any{
					EntryExpression("h.User-Agent"),
					FilterNot("not"),
				},
			},
			want: false,
		},
		{
			name: "sequence of 2 match",
			tks: []any{
				EntryExpression("h.User-Agent"),
				FilterNot("not"),
				FilterTerminator("map"),
			},
			tokenTypes: []any{
				EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
			want: true,
		},
		{
			name: "sequence of don't match",
			tks: []any{
				EntryExpression("h.User-Agent"),
				FilterTerminator("map"),
				FilterNot("not"),
			},
			tokenTypes: []any{
				EntryExpression("h.User-Agent"),
				FilterNot("not"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterTokenMatcher(tt.tks, tt.tokenTypes...); got != tt.want {
				t.Errorf("FilterTokenMatcher() = %v, want %v", got, tt.want)
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
				Not: &EqualOp{LHS: EntryExpression("h.user-agent"), RHS: ConstantExpression("Kmail")},
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
			if diff := cmp.Diff(got, tt.expectedExpression); diff != "" {
				t.Errorf("ParseFilter() expectedExpression %s", diff)
			}
			if diff := cmp.Diff(got1, tt.remaining); diff != "" {
				t.Errorf("ParseFilter() remaining %s", diff)
			}
		})
	}
}
