package funcs

import (
	"os"
	"pimtrace"
	"pimtrace/dataformats/icaldata"
	"testing"
	"time"
)

func TestIcalTemporal(t *testing.T) {
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

	if len(stream) < 2 {
		t.Fatalf("Failed to find VEVENT rows")
	}

	// 2023-10-27T10:30:00 America/New_York (EDT, UTC-4) -> 14:30 UTC
	eventEDT := stream[0]
	dtstartEDT, err := eventEDT.Get("p.DTSTART")
	if err != nil {
		t.Fatalf("Failed to get DTSTART: %v", err)
	}

	expectedEDT := time.Date(2023, 10, 27, 10, 30, 0, 0, time.FixedZone("EDT", -4*3600))
	tEDT, err := temporalCoerce("test", eventEDT, []ValueExpression{
		mockValueExpression{val: dtstartEDT},
	}, nil)
	if err != nil {
		t.Fatalf("temporalCoerce failed: %v", err)
	}
	if tEDT == nil {
		t.Fatalf("temporalCoerce returned nil")
	}
	if !tEDT.Equal(expectedEDT) {
		t.Errorf("temporalCoerce returned %v, expected %v", tEDT, expectedEDT)
	}

	// 2023-11-10T10:30:00 America/New_York (EST, UTC-5) -> 15:30 UTC
	eventEST := stream[1]
	dtstartEST, err := eventEST.Get("p.DTSTART")
	if err != nil {
		t.Fatalf("Failed to get DTSTART: %v", err)
	}

	expectedEST := time.Date(2023, 11, 10, 10, 30, 0, 0, time.FixedZone("EST", -5*3600))
	tEST, err := temporalCoerce("test", eventEST, []ValueExpression{
		mockValueExpression{val: dtstartEST},
	}, nil)
	if err != nil {
		t.Fatalf("temporalCoerce failed: %v", err)
	}
	if tEST == nil {
		t.Fatalf("temporalCoerce returned nil")
	}
	if !tEST.Equal(expectedEST) {
		t.Errorf("temporalCoerce returned %v, expected %v", tEST, expectedEST)
	}

	// Check hour on EDT
	hr := Hour[ValueExpression]{}
	res, err := hr.Run(eventEDT, []ValueExpression{
		mockValueExpression{val: dtstartEDT},
	}, nil)
	if err != nil {
		t.Fatalf("hour function failed: %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok {
		t.Errorf("hour returned %v, expected SimpleIntegerValue", res)
	} else if int(v) != 10 {
		t.Errorf("hour returned %v, expected 10", v)
	}
}
