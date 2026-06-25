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
		Header: map[string]int{"SUMMARY": 0, "DTSTART": 1},
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
		Header: map[string]int{"SUMMARY": 0, "DTSTART": 1},
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
		Header: map[string]int{"SUMMARY": 0},
		ComponentBase: cb,
		Component: comp,
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

	// Test read error / invalid ical? golang-ical is fairly robust, but we can pass invalid data
	badData := `BEGIN:VCALENDAR` // missing END
	rBad := strings.NewReader(badData)
	// golang-ical parses line by line and might just return an empty calendar or error
	_, _ = ReadICalStream(rBad, "ical", "bad.ics") // Just hitting it for coverage
}
