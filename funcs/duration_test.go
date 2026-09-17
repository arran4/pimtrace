package funcs

import (
	"errors"
	"github.com/arran4/golang-ical"
	"pimtrace/dataformats/icaldata"
	"strings"
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func TestDuration(t *testing.T) {
	for _, tc := range []struct {
		name    string
		icsStr  string
		want    *int
		wantErr error
	}{
		{
			name: "UTC timed event DTSTART + DTEND",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART:20231027T100000Z
DTEND:20231027T113000Z
END:VEVENT
END:VCALENDAR`,
			want: intPtr(5400),
		},
		{
			name: "TZ timed event DTSTART + DTEND",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VTIMEZONE
TZID:America/New_York
BEGIN:STANDARD
DTSTART:19701101T020000
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
TZOFFSETFROM:-0400
TZOFFSETTO:-0500
TZNAME:EST
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:19700308T020000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU
TZOFFSETFROM:-0500
TZOFFSETTO:-0400
TZNAME:EDT
END:DAYLIGHT
END:VTIMEZONE
BEGIN:VEVENT
UID:1
DTSTART;TZID=America/New_York:20231027T100000
DTEND;TZID=America/New_York:20231027T113000
END:VEVENT
END:VCALENDAR`,
			want: intPtr(5400),
		},
		{
			name: "Event crossing DST boundary (Fall back)",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VTIMEZONE
TZID:America/New_York
BEGIN:STANDARD
DTSTART:19701101T020000
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
TZOFFSETFROM:-0400
TZOFFSETTO:-0500
TZNAME:EST
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:19700308T020000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU
TZOFFSETFROM:-0500
TZOFFSETTO:-0400
TZNAME:EDT
END:DAYLIGHT
END:VTIMEZONE
BEGIN:VEVENT
UID:1
DTSTART;TZID=America/New_York:20231105T010000
DTEND;TZID=America/New_York:20231105T020000
END:VEVENT
END:VCALENDAR`,
			want: intPtr(7200),
		},
		{
			name: "One-day all-day event",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART;VALUE=DATE:20231027
DTEND;VALUE=DATE:20231028
END:VEVENT
END:VCALENDAR`,
			want: intPtr(86400),
		},
		{
			name: "Multi-day all-day event",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART;VALUE=DATE:20231027
DTEND;VALUE=DATE:20231029
END:VEVENT
END:VCALENDAR`,
			want: intPtr(172800),
		},
		{
			name: "Missing end/duration",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART:20231027T100000Z
END:VEVENT
END:VCALENDAR`,
			wantErr: errors.New("duration: missing duration information"),
		},
		{
			name: "Contradictory DTEND + DURATION",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART:20231027T100000Z
DTEND:20231027T113000Z
DURATION:PT1H30M
END:VEVENT
END:VCALENDAR`,
			wantErr: errors.New("duration: contradictory properties, both DURATION and DTEND present"),
		},
		{
			name: "Explicit DURATION (unsupported by golang-ical)",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART:20231027T100000Z
DURATION:PT1H30M
END:VEVENT
END:VCALENDAR`,
			wantErr: errors.New("duration: golang-ical does not expose a safe public way to interpret a DURATION property: PT1H30M"),
		},
		{
			name: "VTODO with DUE",
			icsStr: `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VTODO
UID:1
DTSTART:20231027T100000Z
DUE:20231027T113000Z
END:VTODO
END:VCALENDAR`,
			want: intPtr(5400),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cal, err := ics.ParseCalendar(strings.NewReader(tc.icsStr))
			if err != nil {
				t.Fatalf("ParseCalendar error: %v", err)
			}

			var comp ics.Component
			var compBase *ics.ComponentBase

			if len(cal.Events()) > 0 {
				comp = cal.Events()[0]
				compBase = &cal.Events()[0].ComponentBase
			} else if len(cal.Todos()) > 0 {
				comp = cal.Todos()[0]
				compBase = &cal.Todos()[0].ComponentBase
			} else {
				t.Fatalf("No component found in fixture")
			}

			icw := &icaldata.ICalWithSource{
				Component:     comp,
				ComponentBase: compBase,
			}

			durFunc := Duration[ValueExpression]{}
			got, err := durFunc.Run(icw, nil, nil)

			if (err != nil) != (tc.wantErr != nil) {
				t.Fatalf("Duration.Run() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && tc.wantErr != nil && err.Error() != tc.wantErr.Error() {
				t.Fatalf("Duration.Run() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err == nil {
				if got == nil && tc.want != nil {
					t.Fatalf("Duration.Run() got nil, want %d", *tc.want)
				}
				if got != nil && tc.want == nil {
					t.Fatalf("Duration.Run() got %v, want nil", got)
				}
				if got != nil && tc.want != nil {
					val := got.Integer()
					if val == nil {
						t.Fatalf("Duration.Run() returned a non-integer value")
					}
					if *val != *tc.want {
						t.Errorf("Duration.Run() = %v, want %v", *val, *tc.want)
					}
				}
			}
		})
	}
}

// This ensures evaluator path integration works
func TestDurationEvaluator(t *testing.T) {
	icsStr := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART:20231027T100000Z
DTEND:20231027T113000Z
END:VEVENT
END:VCALENDAR`
	cal, _ := ics.ParseCalendar(strings.NewReader(icsStr))
	icw := &icaldata.ICalWithSource{
		Component:     cal.Events()[0],
		ComponentBase: &cal.Events()[0].ComponentBase,
	}

	durFunc := Duration[ValueExpression]{}
	val, err := durFunc.Run(icw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	intVal := val.Integer()
	if intVal == nil || *intVal != 5400 {
		t.Fatalf("unexpected output: %v", val)
	}
}

func TestDurationMalformed(t *testing.T) {
	icsStr := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART:20231027T100000Z
DURATION:NOT-A-DURATION
END:VEVENT
END:VCALENDAR`
	cal, _ := ics.ParseCalendar(strings.NewReader(icsStr))
	icw := &icaldata.ICalWithSource{
		Component:     cal.Events()[0],
		ComponentBase: &cal.Events()[0].ComponentBase,
	}

	durFunc := Duration[ValueExpression]{}
	_, err := durFunc.Run(icw, nil, nil)
	if err == nil {
		t.Fatalf("expected error for malformed DURATION but got none")
	}
	if !strings.Contains(err.Error(), "golang-ical does not expose a safe public way to interpret a DURATION property") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestDurationInvalidTimezone(t *testing.T) {
	icsStr := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:1
DTSTART;TZID=Invalid/Timezone:20231027T100000
DTEND;TZID=Invalid/Timezone:20231027T113000
END:VEVENT
END:VCALENDAR`
	cal, _ := ics.ParseCalendar(strings.NewReader(icsStr))
	icw := &icaldata.ICalWithSource{
		Component:     cal.Events()[0],
		ComponentBase: &cal.Events()[0].ComponentBase,
	}

	durFunc := Duration[ValueExpression]{}
	_, err := durFunc.Run(icw, nil, nil)
	if err == nil {
		t.Fatalf("expected error for invalid timezone but got none")
	}
}
