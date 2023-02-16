package funcs

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/goodsign/monday"
	"pimtrace"
	"time"
	"unicode"
)

type Year[T ValueExpression] struct{}

var _ Function[ValueExpression] = Year[ValueExpression]{}

func (c Year[T]) Name() string {
	return "year"
}

func (c Year[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Args:        []Argument{String},
			Description: "Converts time string to a date and returns the year number of that date",
		},
		{
			Args:        []Argument{Integer},
			Description: "Converts Unix time to a date and returns the year number of that date",
		},
	}
}

func (c Year[T]) Run(d pimtrace.Entry, args []T) (pimtrace.Value, error) {
	t, err := Arg1OnlyToTime("year", d, args)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return &pimtrace.SimpleNilValue{}, nil
	}
	return pimtrace.SimpleIntegerValue(int(t.Year())), nil
}

func Arg1OnlyToTime[T ValueExpression](funcName string, d pimtrace.Entry, args []T) (*time.Time, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("%w", ErrExpecting1ArgumentOfTypeStringIntOrDate)
	}
	v, err := args[0].Execute(d)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	if v == nil {
		return nil, fmt.Errorf("%s: %w", funcName, ErrEmptyType)
	}
	var t time.Time
	switch v.(type) {
	case pimtrace.SimpleIntegerValue:
		i := v.Integer()
		if i == nil {
			return nil, fmt.Errorf("%s parse: %w", funcName, ErrNumberError)
		}
		t = time.Unix(int64(*i), 0)
	case pimtrace.SimpleStringValue:
		s := v.String()
		if s == "" {
			return nil, nil
		}
		runes := []rune(s)
		for i, v := range runes {
			if unicode.IsSymbol(v) {
				s = string(runes[:i])
				break
			}
		}
		layout, err := dateparse.ParseFormat(s)
		if err != nil {
			return nil, fmt.Errorf("%s parse format: %w", funcName, err)
		}
		end := len(s)
		if len(layout) < end {
			end = len(layout)
		}
		t, err = monday.NewLocaleDetector().Parse(layout, s[:end])
		if err != nil {
			return nil, fmt.Errorf("%s parse time with locale detector: %w", funcName, err)
		}
	default:
		return nil, fmt.Errorf("%s %w: %s", funcName, ErrUnsupportedType, v.Type())
	}
	return &t, nil
}
