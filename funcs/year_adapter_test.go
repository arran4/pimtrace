package funcs

import (
	"os"
	"pimtrace"
	"pimtrace/dataformats/icaldata"
	"testing"
)

func TestYearAdapter_ICalTimeValue(t *testing.T) {
	f, err := os.Open("testdata/timezone.ics")
	if err != nil {
		t.Fatalf("Failed to open fixture: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	stream, err := icaldata.ReadICalStream(f, "test", "test.ics")
	if err != nil {
		t.Fatalf("Failed to read ICAL stream: %v", err)
	}

	if len(stream) < 3 {
		t.Fatalf("Failed to find New Year's Eve VEVENT row")
	}

	eventNYE := stream[2] // The 3rd event
	val, err := eventNYE.Get("p.DTSTART")
	if err != nil {
		t.Fatalf("Failed to get NYE DTSTART: %v", err)
	}

	// Local time is Dec 31, 2023 at 22:30.
	// UTC time is Jan 1, 2024 at 03:30.
	// YearAdapter should preserve local calendar semantics: Year=2023, Month=12.

	ya := &YearAdapter{}
	resYear, err := ya.Call(val)
	if err != nil {
		t.Fatalf("YearAdapter Call failed: %v", err)
	}

	if y, ok := resYear.(pimtrace.SimpleIntegerValue); !ok {
		t.Errorf("YearAdapter returned non-integer: %v", resYear)
	} else if int(y) != 2023 {
		t.Errorf("YearAdapter returned %v, expected 2023", y)
	}

	ma := &MonthAdapter{}
	resMonth, err := ma.Call(val)
	if err != nil {
		t.Fatalf("MonthAdapter Call failed: %v", err)
	}

	if m, ok := resMonth.(pimtrace.SimpleIntegerValue); !ok {
		t.Errorf("MonthAdapter returned non-integer: %v", resMonth)
	} else if int(m) != 12 {
		t.Errorf("MonthAdapter returned %v, expected 12", m)
	}
}
