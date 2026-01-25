package funcs

import (
	"fmt"
	"pimtrace"
	"time"

	"github.com/araddon/dateparse"
	"github.com/goodsign/monday"
)

type MonthAdapter struct{}

func (m *MonthAdapter) Call(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected 1 argument")
	}

	arg := args[0]
	if arg == nil {
		return nil, nil
	}

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
		// Simplified parsing compared to YearAdapter for brevity/directness,
		// but should ideally match logic. Resusing common logic would be better.
		// For now implementing basic parsing.
		t, err = dateparse.ParseAny(v)
		if err != nil {
			// Try monday
			var layout string
			layout, err = dateparse.ParseFormat(v)
			if err == nil {
				t, err = monday.NewLocaleDetector().Parse(layout, v)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
	case pimtrace.Value:
		if i := v.Integer(); i != nil {
			t = time.Unix(int64(*i), 0)
		} else {
			return m.Call(v.String())
		}
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}

	return pimtrace.SimpleIntegerValue(int(t.Month())), nil
}
