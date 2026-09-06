package icaldata

import (
	"pimtrace"
	"reflect"
	"strings"
	"testing"

	ics "github.com/arran4/golang-ical"
)

func TestICalWithSourceGetSize(t *testing.T) {
	cb := &ics.ComponentBase{
		Properties: []ics.IANAProperty{
			{BaseProperty: ics.BaseProperty{IANAToken: "SUMMARY"}},
			{BaseProperty: ics.BaseProperty{IANAToken: "LOCATION"}},
		},
	}
	ic := &ICalWithSource{ComponentBase: cb, Header: map[string]int{"SUMMARY": 0, "LOCATION": 1}}
	v, err := ic.Get("sz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i := v.Integer()
	if i == nil || *i != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}

func TestICalWithSource_Self(t *testing.T) {
	r := &ICalWithSource{}
	if r.Self() != r {
		t.Errorf("Self() didn't return same pointer")
	}
}

func TestICalWithSource_HeadersStringArray(t *testing.T) {
	r := &ICalWithSource{
		Header: map[string]int{"SUMMARY": 0, "DTSTART": 1},
	}
	res := r.HeadersStringArray()
	if len(res) != 2 {
		t.Errorf("HeadersStringArray() returned wrong length")
	}
}

func TestICalWithSource_StringArray(t *testing.T) {
	cb := &ics.ComponentBase{
		Properties: []ics.IANAProperty{
			{BaseProperty: ics.BaseProperty{IANAToken: string(ics.PropertySummary), Value: "Meeting"}},
			{BaseProperty: ics.BaseProperty{IANAToken: string(ics.PropertyDtstart), Value: "20231027T100000Z"}},
		},
	}

	r := &ICalWithSource{
		Header:        map[string]int{"SUMMARY": 0, "DTSTART": 1},
		ComponentBase: cb,
	}

	res := r.StringArray([]string{"SUMMARY", "NONEXISTENT", "DTSTART"})
	expected := []string{"Meeting", "20231027T100000Z"}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("StringArray() = %v, want %v", res, expected)
	}
}

func TestICalWithSource_Get(t *testing.T) {
	cb := &ics.ComponentBase{
		Properties: []ics.IANAProperty{
			{BaseProperty: ics.BaseProperty{IANAToken: string(ics.PropertySummary), Value: "Meeting"}},
			{BaseProperty: ics.BaseProperty{IANAToken: string(ics.PropertyDtstart), Value: "20231027T100000Z"}},
		},
	}

	r := &ICalWithSource{
		Header:        map[string]int{"SUMMARY": 0, "DTSTART": 1},
		ComponentBase: cb,
	}

	// sz
	v, err := r.Get("sz")
	if err != nil {
		t.Errorf("Get(sz) error: %v", err)
	}
	if iv, ok := v.(pimtrace.SimpleIntegerValue); !ok || int(iv) != 2 {
		t.Errorf("Get(sz) = %v, want 2", v)
	}

	// property
	v, err = r.Get("SUMMARY.val")
	if err != nil {
		t.Errorf("Get(SUMMARY.val) error: %v", err)
	}
	if sv, ok := v.(pimtrace.SimpleStringValue); !ok || string(sv) != "Meeting" {
		t.Errorf("Get(SUMMARY.val) = %v, want Meeting", v)
	}

	// missing property
	_, err = r.Get("NONEXISTENT.val")
	if err == nil {
		t.Errorf("Get(NONEXISTENT.val) expected error")
	}

	// short key
	_, err = r.Get("SUMMARY")
	if err == nil {
		t.Errorf("Get(SUMMARY) expected error (too short)")
	}
}

func TestData(t *testing.T) {
	var d Data = make([]*ICalWithSource, 0)

	if d.Len() != 0 {
		t.Errorf("Len() = %v, want 0", d.Len())
	}

	r1 := &ICalWithSource{}
	r2 := &ICalWithSource{}

	d = d.SetEntry(0, r1).(Data)
	if d.Len() != 1 {
		t.Errorf("Len() = %v, want 1", d.Len())
	}

	d = d.SetEntry(2, r2).(Data) // Should pad
	if d.Len() != 3 {
		t.Errorf("Len() = %v, want 3", d.Len())
	}
	if d.Entry(2) != r2 {
		t.Errorf("Entry(2) != r2")
	}
	if d.Entry(1) != (*ICalWithSource)(nil) {
		t.Errorf("Entry(1) should be nil padded")
	}
	if d.Entry(5) != nil {
		t.Errorf("Entry(5) should be nil (out of bounds)")
	}

	d = d.Truncate(2).(Data)
	if d.Len() != 2 {
		t.Errorf("Truncate(2) Len = %v, want 2", d.Len())
	}

	if d.Self() == nil {
		t.Errorf("Self() returned nil")
	}

	dNew := d.NewSelf()
	if dNew.Len() != 0 {
		t.Errorf("NewSelf() Len = %v, want 0", dNew.Len())
	}
}

func TestData_Output(t *testing.T) {
	cb := &ics.ComponentBase{
		Properties: []ics.IANAProperty{
			{BaseProperty: ics.BaseProperty{IANAToken: string(ics.PropertySummary), Value: "Meeting"}},
		},
	}
	comp := &ics.VEvent{ComponentBase: *cb}

	var d Data = make([]*ICalWithSource, 0)
	d = append(d, &ICalWithSource{
		Header:        map[string]int{"SUMMARY": 0},
		ComponentBase: cb,
		Component:     comp,
	})

	// Test CSV and Table streams (writing to - stdout)
	//d.WriteCSVFile("-")
	//d.WriteTableFile("-")

	// Test ICal stream
	err := d.WriteICalFile("test.ics")
	if err != nil {
		t.Errorf("WriteICalFile error: %v", err)
	}
}

func TestReadICalStream(t *testing.T) {
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Example Corp.//Cal//EN
BEGIN:VEVENT
UID:12345
DTSTAMP:20231027T100000Z
SUMMARY:Test Event
END:VEVENT
END:VCALENDAR`

	r := strings.NewReader(icalData)
	sources, err := ReadICalStream(r, "ical", "test.ics")
	if err != nil {
		t.Errorf("ReadICalStream error: %v", err)
	}
	if len(sources) != 1 {
		t.Errorf("ReadICalStream expected 1 event, got %d", len(sources))
	}
	if sources[0].SourceType != "ical" || sources[0].SourceFile != "test.ics" {
		t.Errorf("ReadICalStream source meta incorrect")
	}

	// Test regression fixture: VTIMEZONE + VEVENT to ensure it doesn't panic
	// and timezone is ignored while event is kept, while retaining timestamp interpretation.
	tzData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Example Corp.//Cal//EN
BEGIN:VTIMEZONE
TZID:America/New_York
BEGIN:STANDARD
DTSTART:20071104T020000
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
TZOFFSETFROM:-0400
TZOFFSETTO:-0500
TZNAME:EST
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
UID:12345
DTSTAMP:20231027T100000Z
DTSTART;TZID=America/New_York:20231027T100000
SUMMARY:Test Event
END:VEVENT
END:VCALENDAR`

	rTz := strings.NewReader(tzData)
	tzSources, tzErr := ReadICalStream(rTz, "ical", "tz.ics")
	if tzErr != nil {
		t.Errorf("ReadICalStream (regression) error: %v", tzErr)
	}
	if len(tzSources) != 1 {
		t.Errorf("ReadICalStream (regression) expected exactly 1 event (skipping VTIMEZONE), got %d", len(tzSources))
	} else if ev, ok := tzSources[0].Component.(*ics.VEvent); !ok {
		t.Errorf("ReadICalStream (regression) expected VEVENT, got %T", tzSources[0].Component)
	} else {
		startTime, err := ev.GetStartAt()
		if err != nil {
			t.Errorf("ReadICalStream (regression) GetStartAt error: %v", err)
		} else if startTime.UTC().Format("2006-01-02 15:04:05 -0700 MST") != "2023-10-27 14:00:00 +0000 UTC" {
			t.Errorf("ReadICalStream (regression) GetStartAt expected '2023-10-27 14:00:00 +0000 UTC', got %v", startTime.UTC().Format("2006-01-02 15:04:05 -0700 MST"))
		}
	}

	// Test coverage for another supported component (VTODO)
	todoData := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VTODO
UID:67890
SUMMARY:Task
END:VTODO
END:VCALENDAR`

	rTodo := strings.NewReader(todoData)
	todoSources, todoErr := ReadICalStream(rTodo, "ical", "todo.ics")
	if todoErr != nil {
		t.Errorf("ReadICalStream (VTODO) error: %v", todoErr)
	}
	if len(todoSources) != 1 {
		t.Errorf("ReadICalStream (VTODO) expected 1 event, got %d", len(todoSources))
	}

	// Test handling of an unsupported component type (e.g. unknown component)
	// which returns a controlled error rather than panicking.
	unsupportedData := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:X-UNKNOWN
UID:99999
END:X-UNKNOWN
END:VCALENDAR`
	rUnknown := strings.NewReader(unsupportedData)
	_, unknownErr := ReadICalStream(rUnknown, "ical", "unknown.ics")
	if unknownErr == nil {
		t.Errorf("ReadICalStream (unknown component) expected error, got nil")
	} else if !strings.Contains(unknownErr.Error(), "unsupported component type") {
		t.Errorf("ReadICalStream (unknown component) expected 'unsupported component type' error, got %v", unknownErr)
	}

	// Test read error / invalid ical? golang-ical is fairly robust, but we can pass invalid data
	badData := `BEGIN:VCALENDAR` // missing END
	rBad := strings.NewReader(badData)
	// golang-ical parses line by line and might just return an empty calendar or error
	_, _ = ReadICalStream(rBad, "ical", "bad.ics") // Just hitting it for coverage
}
