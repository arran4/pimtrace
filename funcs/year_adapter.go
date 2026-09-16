package funcs

import (
	"fmt"
	"pimtrace"
	"time"
)

type YearAdapter struct{}

func (y *YearAdapter) Call(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected 1 argument")
	}

	arg := args[0]
	if arg == nil {
		return nil, nil
	}

	var t *time.Time
	var err error

	switch v := arg.(type) {
	case int:
		vt := time.Unix(int64(v), 0)
		t = &vt
	case int64:
		vt := time.Unix(v, 0)
		t = &vt
	case string:
		if v == "" {
			return nil, fmt.Errorf("empty string")
		}
		t, err = pimtrace.ParseDate(v)
		if err != nil {
			return nil, err
		}
	case pimtrace.SimpleNilValue:
		return nil, nil
	case *pimtrace.SimpleNilValue:
		return nil, nil
	case pimtrace.SimpleIntegerValue:
		i := v.Integer()
		vt := time.Unix(int64(*i), 0)
		t = &vt
	case pimtrace.SimpleStringValue:
		if v.String() == "" {
			return nil, fmt.Errorf("empty string")
		}
		t, err = pimtrace.ParseDate(v.String())
		if err != nil {
			return nil, err
		}
	case pimtrace.Value:
		if v.Type() == pimtrace.Integer {
			i := v.Integer()
			if i != nil {
				vt := time.Unix(int64(*i), 0)
				t = &vt
			}
		} else if v.Type() != pimtrace.Float && v.Time() != nil {
			t = v.Time()
		} else if v.Type() == pimtrace.String {
			if v.String() == "" {
				return nil, fmt.Errorf("empty string")
			}
			t, err = pimtrace.ParseDate(v.String())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported type %T", v)
		}
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}

	if t == nil {
		return nil, nil
	}
	return pimtrace.SimpleIntegerValue(int(t.Year())), nil
}
