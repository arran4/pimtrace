package ast

import (
	"pimtrace"
	"testing"
)

func TestComparatorAdapter(t *testing.T) {
	tests := []struct {
		name    string
		a       pimtrace.Value
		b       pimtrace.Value
		want    int
		wantErr bool
	}{
		{
			name:    "numeric: 2 < 10",
			a:       pimtrace.SimpleStringValue("2"),
			b:       pimtrace.SimpleStringValue("10"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "numeric: 10 > 2",
			a:       pimtrace.SimpleStringValue("10"),
			b:       pimtrace.SimpleStringValue("2"),
			want:    1,
			wantErr: false,
		},
		{
			name:    "numeric: 2.5 < 10.25",
			a:       pimtrace.SimpleStringValue("2.5"),
			b:       pimtrace.SimpleStringValue("10.25"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "numeric: negatives -5 < 2",
			a:       pimtrace.SimpleStringValue("-5"),
			b:       pimtrace.SimpleStringValue("2"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "numeric: zero 0 == 0",
			a:       pimtrace.SimpleStringValue("0"),
			b:       pimtrace.SimpleStringValue("0.0"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "string lexical: a < b",
			a:       pimtrace.SimpleStringValue("a"),
			b:       pimtrace.SimpleStringValue("b"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "numeric with whitespace",
			a:       pimtrace.SimpleStringValue(" 5 "),
			b:       pimtrace.SimpleStringValue(" 10 "),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "mixed coercion error",
			a:       pimtrace.SimpleStringValue("10"),
			b:       pimtrace.SimpleStringValue("abc"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string against numeric is error",
			a:       pimtrace.SimpleStringValue("10"),
			b:       pimtrace.SimpleStringValue(""),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ComparatorAdapter{Value: tt.a}
			got, err := c.Compare(ComparatorAdapter{Value: tt.b})
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparatorAdapterDates(t *testing.T) {
	tests := []struct {
		name    string
		a       pimtrace.Value
		b       pimtrace.Value
		want    int
		wantErr bool
	}{
		{
			name:    "date: early < late",
			a:       pimtrace.SimpleStringValue("2020-01-01"),
			b:       pimtrace.SimpleStringValue("2020-01-02"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "date: late > early",
			a:       pimtrace.SimpleStringValue("2020-01-02"),
			b:       pimtrace.SimpleStringValue("2020-01-01"),
			want:    1,
			wantErr: false,
		},
		{
			name:    "date mixed formats: ISO < RFC3339",
			a:       pimtrace.SimpleStringValue("2020-01-01T00:00:00Z"),
			b:       pimtrace.SimpleStringValue("2020-01-02 00:00:00"),
			want:    -1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ComparatorAdapter{Value: tt.a}
			got, err := c.Compare(ComparatorAdapter{Value: tt.b})
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
