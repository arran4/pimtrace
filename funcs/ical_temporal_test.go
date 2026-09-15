package funcs

import (
	"github.com/arran4/golang-ical"
	"os"
	"pimtrace"
	"pimtrace/dataformats/icaldata"
	"testing"
	"time"
)

func TestIcalTemporal(t *testing.T) {
	// The requirement is that the `.ics` fixture containing VTIMEZONE and a timezone-bearing VEVENT
	// should be exercised through the actual PIMTrace / golang-ical ingestion path.
	// Where golang-ical already provides parsed time semantics, use those APIs rather than reparsing raw iCalendar strings.

	f, err := os.Open("testdata/timezone.ics")
	if err != nil {
		t.Fatalf("Failed to open fixture: %v", err)
	}
	defer f.Close()

	stream, err := icaldata.ReadICalStream(f, "test", "test.ics")
	if err != nil {
		t.Fatalf("Failed to read ICAL stream: %v", err)
	}

	if len(stream) == 0 {
		t.Fatalf("Failed to find VEVENT row")
	}

	event := stream[0]
	ve, ok := event.Component.(*ics.VEvent)
	if !ok {
		t.Fatalf("Component is not a VEvent")
	}

	dtstart, err := ve.GetStartAt()
	if err != nil {
		t.Fatalf("Failed to parse start at: %v", err)
	}

	// 2023-10-27T10:30:00 America/New_York (EDT/EST crossover context, EDT is -0400 on Oct 27)
	// New York is EDT (UTC-4) in Oct 2023. So 10:30 EDT == 14:30 UTC.
	if !dtstart.Equal(time.Date(2023, 10, 27, 14, 30, 0, 0, time.UTC)) {
		t.Errorf("golang-ical GetStartAt() returned %v, expected 2023-10-27 14:30:00 UTC", dtstart)
	}

	// We pass the parsed time to our temporal function through the Unix timestamp representation to verify
	// the shared coercion policy preserves it without modifying its global CoerceDate behavior.
	// This also simulates a valid Value extraction where the time is pre-parsed.
	t1, err := temporalCoerce("test", event, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleIntegerValue(int(dtstart.Unix()))},
	}, nil)

	if err != nil {
		t.Fatalf("temporalCoerce failed: %v", err)
	}

	if t1 == nil {
		t.Fatalf("temporalCoerce returned nil")
	}

	if !t1.Equal(dtstart) {
		t.Errorf("temporalCoerce returned %v, expected %v", t1, dtstart)
	}
}
