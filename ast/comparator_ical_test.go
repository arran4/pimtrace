package ast

import (
	"pimtrace/dataformats/icaldata"
	"testing"
	"os"
	"pimtrace"
)

func TestICalTimeValue_Comparison(t *testing.T) {
	// The PR requested a regression testing that comparing two timezone-bearing DTSTART values through
	// the comparator works and does not fall back to float parsing, plus comparing against a string date literal.

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

	eventEDT := stream[0]
	eventEST := stream[1]

	val1, err := eventEDT.Get("p.DTSTART")
	if err != nil {
		t.Fatalf("Failed to get EDT DTSTART: %v", err)
	}

	val2, err := eventEST.Get("p.DTSTART")
	if err != nil {
		t.Fatalf("Failed to get EST DTSTART: %v", err)
	}

	// Direct Less() / Equal() tests.
	if val1.Equal(val2) {
		t.Errorf("EDT value and EST value should not be equal")
	}

	if !val1.Less(val2) {
		t.Errorf("EDT value (Oct) should be less than EST value (Nov)")
	}
	if val2.Less(val1) {
		t.Errorf("EST value (Nov) should not be less than EDT value (Oct)")
	}

	// Compare against a string date literal wrapped as SimpleStringValue
	literalVal := pimtrace.SimpleStringValue("Wed, 01 Nov 2023 00:00:00 +0000")
	// Oct 27 < Nov 1
	if !val1.Less(literalVal) {
		t.Errorf("EDT value (Oct 27) should be less than Nov 1 literal")
	}
	// Nov 10 > Nov 1
	if val2.Less(literalVal) {
		t.Errorf("EST value (Nov 10) should not be less than Nov 1 literal")
	}

	// Through ComparatorAdapter explicitly
	cmpAdapter1 := ComparatorAdapter{Value: val1}
	cmpAdapter2 := ComparatorAdapter{Value: val2}

	isLess, err := cmpAdapter1.Compare(cmpAdapter2.Value)
	if err != nil {
		t.Fatalf("ComparatorAdapter '<' failed: %v", err)
	}
	if isLess >= 0 {
		t.Errorf("ComparatorAdapter expected val1 < val2")
	}

	isGreater, err := cmpAdapter2.Compare(cmpAdapter1.Value)
	if err != nil {
		t.Fatalf("ComparatorAdapter '>' failed: %v", err)
	}
	if isGreater <= 0 {
		t.Errorf("ComparatorAdapter expected val2 > val1")
	}

	// Test comparison with ordinary literal value (evaluator does interface comparisons via ComparatorAdapter)
	isLessLiteral, err := cmpAdapter1.Compare(literalVal)
	if err != nil {
		t.Fatalf("ComparatorAdapter literal '<' failed: %v", err)
	}
	if isLessLiteral >= 0 {
		t.Errorf("ComparatorAdapter expected val1 < literalVal")
	}
}
