package funcs

import (
	"fmt"
	"pimtrace"
	"time"
	"unicode"

	"github.com/araddon/dateparse"
	"github.com/goodsign/monday"
)

type YearAdapter struct{}

func (y *YearAdapter) Call(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected 1 argument")
	}

	arg := args[0]
	if arg == nil {
		return nil, nil // Or typed nil?
	}

	// pimtrace functions work with pimtrace.Value (mostly).
	// But go-evaluator uses raw interface{}.
	// The arguments passed here are evaluated values.

	// Logic ported from Arg1OnlyToTime but for interface{}
	var t time.Time
	var err error

	switch v := arg.(type) {
	case int:
		t = time.Unix(int64(v), 0)
	case int64:
		t = time.Unix(v, 0)
	case string:
		if v == "" {
			return nil, nil
		}
		s := v
		// Strip symbol logic
		runes := []rune(s)
		for i, r := range runes {
			if unicode.IsSymbol(r) {
				s = string(runes[:i])
				break
			}
		}

		var layout string
		layout, err = dateparse.ParseFormat(s)
		if err != nil {
			return nil, fmt.Errorf("parse format: %w", err)
		}
		end := len(s)
		if len(layout) < end {
			end = len(layout)
		}
		t, err = monday.NewLocaleDetector().Parse(layout, s[:end])
		if err != nil {
			return nil, fmt.Errorf("parse time with locale detector: %w", err)
		}
	case pimtrace.Value: // Handle pimtrace.Value if args arrive as such
		if i := v.Integer(); i != nil {
			t = time.Unix(int64(*i), 0)
		} else {
			// recurse or handle string
			return y.Call(v.String())
		}
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}

	return pimtrace.SimpleIntegerValue(int(t.Year())), nil
}
