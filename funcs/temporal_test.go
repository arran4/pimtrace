package funcs

import (
	"errors"
	"pimtrace"
	"testing"
	"time"
)

func TestTemporalCoerce(t *testing.T) {
	for _, tc := range []struct {
		name    string
		args    []ValueExpression
		want    *time.Time
		wantErr error
	}{
		{
			name: "RFC3339 UTC",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27T10:00:00Z")},
			},
			want: timePtr(time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)),
		},
		{
			name: "Timestamp with numeric timezone offset",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27T10:00:00+02:00")},
			},
			want: timePtr(time.Date(2023, 10, 27, 10, 0, 0, 0, time.FixedZone("", 2*3600))),
		},
		{
			name: "RFC/mail date",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("Wed, 8 Feb 2023 19:00:46 +1100")},
			},
			want: timePtr(time.Date(2023, 2, 8, 19, 0, 46, 0, time.FixedZone("", 11*3600))),
		},
		{
			name: "Existing accepted ordinary string date",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27")},
			},
			want: timePtr(time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "Unix integer timestamp",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleIntegerValue(1675843200)}, // 2023-02-08 08:00:00 UTC
			},
			want: timePtr(time.Unix(1675843200, 0)),
		},
		{
			name: "Compact iCalendar DATE",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("20231027")},
			},
			want: timePtr(time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "Compact iCalendar DATE-TIME",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("20231027T103000Z")},
			},
			want: timePtr(time.Date(2023, 10, 27, 10, 30, 0, 0, time.UTC)),
		},
		{
			name: "Compact iCalendar DATE-TIME without Z",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("20231027T103000")},
			},
			want: timePtr(time.Date(2023, 10, 27, 10, 30, 0, 0, time.UTC)),
		},
		{
			name: "Empty string",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("")},
			},
			wantErr: errors.New("empty date string"),
		},
		{
			name: "Invalid date string",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("invalid-date")},
			},
			wantErr: errors.New("parse date"),
		},
		{
			name: "Nil/missing value",
			args: []ValueExpression{
				mockValueExpression{val: &pimtrace.SimpleNilValue{}},
			},
			want: nil,
		},
		{
			name: "Unsupported value type (Float)",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleFloatValue(123.45)},
			},
			wantErr: ErrUnsupportedType,
		},
		{
			name: "Expression evaluation error propagation",
			args: []ValueExpression{
				mockValueExpression{err: errors.New("eval err")},
			},
			wantErr: errors.New("eval err"),
		},
		{
			name: "DST/timezone-sensitive case (AEDT -> UTC)",
			args: []ValueExpression{
				mockValueExpression{val: pimtrace.SimpleStringValue("Wed, 8 Feb 2023 19:00:46 +1100 (AEDT)")},
			},
			want: timePtr(time.Date(2023, 2, 8, 19, 0, 46, 0, time.FixedZone("", 11*3600))),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := temporalCoerce("testfunc", nil, tc.args, nil)
			if (err != nil) != (tc.wantErr != nil) {
				t.Errorf("temporalCoerce() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil {
				if got == nil && tc.want != nil {
					t.Errorf("temporalCoerce() got nil, want %v", tc.want)
				} else if got != nil && tc.want == nil {
					t.Errorf("temporalCoerce() got %v, want nil", got)
				} else if got != nil && tc.want != nil {
					if !got.Equal(*tc.want) {
						t.Errorf("temporalCoerce() got %v, want %v", got, tc.want)
					}
				}
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
