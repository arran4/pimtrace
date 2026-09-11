package ast

import (
	"pimtrace"
	"testing"
)

func TestComparatorAdapter(t *testing.T) {
	tests := []struct {
		name    string
		a       pimtrace.Value
		b       any
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
		{
			name:    "float64 input without integer truncation: 10.25 < 10.5",
			a:       pimtrace.SimpleStringValue("10.25"),
			b:       float64(10.5),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "float64 input without integer truncation: 10.5 > 10.25",
			a:       pimtrace.SimpleStringValue("10.5"),
			b:       float64(10.25),
			want:    1,
			wantErr: false,
		},
		{
			name:    "SimpleFloatValue vs SimpleFloatValue: 10.25 < 10.5",
			a:       pimtrace.SimpleFloatValue(10.25),
			b:       pimtrace.SimpleFloatValue(10.5),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "SimpleFloatValue vs float64: 10.25 == 10.25",
			a:       pimtrace.SimpleFloatValue(10.25),
			b:       float64(10.25),
			want:    0,
			wantErr: false,
		},
		{
			name:    "SimpleFloatValue vs SimpleIntegerValue: 2.5 > 2",
			a:       pimtrace.SimpleFloatValue(2.5),
			b:       pimtrace.SimpleIntegerValue(2),
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ComparatorAdapter{Value: tt.a}
			got, err := c.Compare(tt.b)
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
		{
			name:    "valid date vs invalid text -> error",
			a:       pimtrace.SimpleStringValue("2020-01-01"),
			b:       pimtrace.SimpleStringValue("apple"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid text vs valid date -> error",
			a:       pimtrace.SimpleStringValue("apple"),
			b:       pimtrace.SimpleStringValue("2020-01-01"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty vs date -> error",
			a:       pimtrace.SimpleStringValue(""),
			b:       pimtrace.SimpleStringValue("2020-01-01"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "date vs empty -> error",
			a:       pimtrace.SimpleStringValue("2020-01-01"),
			b:       pimtrace.SimpleStringValue(""),
			want:    0,
			wantErr: true,
		},
		{
			name:    "all-day compact DATE vs ISO date: 20200102 < 2020-01-03",
			a:       pimtrace.SimpleStringValue("20200102"),
			b:       pimtrace.SimpleStringValue("2020-01-03"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "ISO date vs all-day compact DATE: 2020-01-01 < 20200102",
			a:       pimtrace.SimpleStringValue("2020-01-01"),
			b:       pimtrace.SimpleStringValue("20200102"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "all-day compact DATE vs RFC3339: 20200102 < 2020-01-02T15:04:05Z",
			a:       pimtrace.SimpleStringValue("20200102"),
			b:       pimtrace.SimpleStringValue("2020-01-02T15:04:05Z"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "RFC3339 vs all-day compact DATE: 2020-01-02T15:04:05Z > 20200102",
			a:       pimtrace.SimpleStringValue("2020-01-02T15:04:05Z"),
			b:       pimtrace.SimpleStringValue("20200102"),
			want:    1,
			wantErr: false,
		},
		{
			name:    "compact DATE vs compact DATE-TIME: 20200102 < 20200102T090000Z",
			a:       pimtrace.SimpleStringValue("20200102"),
			b:       pimtrace.SimpleStringValue("20200102T090000Z"),
			want:    -1,
			wantErr: false,
		},
		{
			name:    "compact DATE vs ISO date equal: 20200102 == 2020-01-02",
			a:       pimtrace.SimpleStringValue("20200102"),
			b:       pimtrace.SimpleStringValue("2020-01-02"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "ISO date vs compact DATE equal: 2020-01-02 == 20200102",
			a:       pimtrace.SimpleStringValue("2020-01-02"),
			b:       pimtrace.SimpleStringValue("20200102"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "compact DATE against itself: 20200102 == 20200102",
			a:       pimtrace.SimpleStringValue("20200102"),
			b:       pimtrace.SimpleStringValue("20200102"),
			want:    0,
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

func TestDetermineComparisonMode(t *testing.T) {
	tests := []struct {
		name string
		a    pimtrace.Value
		b    pimtrace.Value
		want comparisonMode
	}{
		{
			name: "compact DATE vs ISO date -> modeDate",
			a:    pimtrace.SimpleStringValue("20200102"),
			b:    pimtrace.SimpleStringValue("2020-01-03"),
			want: modeDate,
		},
		{
			name: "ISO date vs compact DATE -> modeDate",
			a:    pimtrace.SimpleStringValue("2020-01-01"),
			b:    pimtrace.SimpleStringValue("20200102"),
			want: modeDate,
		},
		{
			name: "compact DATE vs RFC3339 -> modeDate",
			a:    pimtrace.SimpleStringValue("20200102"),
			b:    pimtrace.SimpleStringValue("2020-01-02T15:04:05Z"),
			want: modeDate,
		},
		{
			name: "RFC3339 vs compact DATE -> modeDate",
			a:    pimtrace.SimpleStringValue("2020-01-02T15:04:05Z"),
			b:    pimtrace.SimpleStringValue("20200102"),
			want: modeDate,
		},
		{
			name: "ISO date vs ISO date -> modeDate",
			a:    pimtrace.SimpleStringValue("2020-01-01"),
			b:    pimtrace.SimpleStringValue("2020-01-02"),
			want: modeDate,
		},
		{
			name: "numeric strings 2 vs 10 -> modeNumeric",
			a:    pimtrace.SimpleStringValue("2"),
			b:    pimtrace.SimpleStringValue("10"),
			want: modeNumeric,
		},
		{
			name: "decimals 2.5 vs 10.25 -> modeNumeric",
			a:    pimtrace.SimpleStringValue("2.5"),
			b:    pimtrace.SimpleStringValue("10.25"),
			want: modeNumeric,
		},
		{
			name: "SimpleFloatValues -> modeNumeric",
			a:    pimtrace.SimpleFloatValue(10.25),
			b:    pimtrace.SimpleFloatValue(10.5),
			want: modeNumeric,
		},
		{
			name: "compact numbers/dates 20200102 vs 20200103 -> modeNumeric",
			a:    pimtrace.SimpleStringValue("20200102"),
			b:    pimtrace.SimpleStringValue("20200103"),
			want: modeNumeric,
		},
		{
			name: "numeric vs invalid text 10 vs abc -> modeNumeric",
			a:    pimtrace.SimpleStringValue("10"),
			b:    pimtrace.SimpleStringValue("abc"),
			want: modeNumeric,
		},
		{
			name: "invalid text vs numeric abc vs 10 -> modeNumeric",
			a:    pimtrace.SimpleStringValue("abc"),
			b:    pimtrace.SimpleStringValue("10"),
			want: modeNumeric,
		},
		{
			name: "ISO date vs invalid text -> modeDate",
			a:    pimtrace.SimpleStringValue("2020-01-01"),
			b:    pimtrace.SimpleStringValue("apple"),
			want: modeDate,
		},
		{
			name: "invalid text vs ISO date -> modeDate",
			a:    pimtrace.SimpleStringValue("apple"),
			b:    pimtrace.SimpleStringValue("2020-01-01"),
			want: modeDate,
		},
		{
			name: "empty vs ISO date -> modeDate",
			a:    pimtrace.SimpleStringValue(""),
			b:    pimtrace.SimpleStringValue("2020-01-01"),
			want: modeDate,
		},
		{
			name: "ISO date vs empty -> modeDate",
			a:    pimtrace.SimpleStringValue("2020-01-01"),
			b:    pimtrace.SimpleStringValue(""),
			want: modeDate,
		},
		{
			name: "lexical strings apple vs banana -> modeLexical",
			a:    pimtrace.SimpleStringValue("apple"),
			b:    pimtrace.SimpleStringValue("banana"),
			want: modeLexical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineComparisonMode(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("determineComparisonMode(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
