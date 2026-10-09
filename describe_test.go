package pimtrace_test

import (
	"bytes"

	"pimtrace/dataformats"
	"pimtrace/dataformats/icaldata"
	"pimtrace/dataformats/maildata"
	"pimtrace/dataformats/tabledata"
	"strings"
	"testing"
)

func TestDescribeCSV(t *testing.T) {
	csvData := "A,B,C.d\n1,2,3\n4,,6\n"
	r := strings.NewReader(csvData)
	rows, err := tabledata.ReadCSV(r, "csv", "-")
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	d := tabledata.Data(rows)

	descs := dataformats.Describe(d, 100)

	expectedKeys := map[string]bool{
		"c.A":   false,
		"c.B":   false,
		"c.C.d": false,
	}

	for _, desc := range descs {
		if _, ok := expectedKeys[desc.Expression]; ok {
			expectedKeys[desc.Expression] = true
		}

		if desc.MatchCount != 2 {
			t.Errorf("Expected count 2 for %s, got %d", desc.Expression, desc.MatchCount)
		}
	}

	for k, found := range expectedKeys {
		if !found {
			t.Errorf("Expected to find field expression %s in describe output", k)
		}
	}
}

func TestDescribeMail(t *testing.T) {
	emlData := "From: alice@example.com\nTo: bob@example.com\nSubject: Secret Subject 123\nX-Custom: hello\n\nBody"
	r := strings.NewReader(emlData)

	msgs, err := maildata.ReadMailStream(r, "mailfile", "-")
	if err != nil {
		t.Fatalf("Failed to read Mail: %v", err)
	}

	d := maildata.Data(msgs)
	descs := dataformats.Describe(d, 100)

	expectedKeys := map[string]bool{
		"h.From":     false,
		"h.To":       false,
		"h.Subject":  false,
		"h.X-Custom": false,
	}

	for _, desc := range descs {
		if _, ok := expectedKeys[desc.Expression]; ok {
			expectedKeys[desc.Expression] = true
		}
	}

	for k, found := range expectedKeys {
		if !found {
			t.Errorf("Expected to find field expression %s in describe output", k)
		}
	}
}

func TestDescribeICal(t *testing.T) {
	icsData := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Example Corp//NONSGML My Product//EN\nBEGIN:VEVENT\nUID:19970901T130000Z-123401@example.com\nDTSTAMP:19970901T130000Z\nDTSTART:19970903T163000Z\nDTEND:19970903T190000Z\nSUMMARY:Secret Meeting\nCLASS:PRIVATE\nX-MY-CUSTOM:Hello\nEND:VEVENT\nEND:VCALENDAR\n"
	r := strings.NewReader(icsData)

	comps, err := icaldata.ReadICalStream(r, "ical", "-")
	if err != nil {
		t.Fatalf("Failed to read ICal: %v", err)
	}

	d := icaldata.Data(comps)
	descs := dataformats.Describe(d, 100)

	expectedKeys := map[string]bool{
		"p.UID":         false,
		"p.DTSTAMP":     false,
		"p.DTSTART":     false,
		"p.DTEND":       false,
		"p.SUMMARY":     false,
		"p.CLASS":       false,
		"p.X-MY-CUSTOM": false,
	}

	for _, desc := range descs {
		if _, ok := expectedKeys[desc.Expression]; ok {
			expectedKeys[desc.Expression] = true
		}
	}

	for k, found := range expectedKeys {
		if !found {
			t.Errorf("Expected to find field expression %s in describe output", k)
		}
	}
}

func TestDescribePrivacy(t *testing.T) {
	csvData := "A,B\nsecret_value_1,secret_value_2\n"
	r := strings.NewReader(csvData)
	rows, _ := tabledata.ReadCSV(r, "csv", "-")
	d := tabledata.Data(rows)

	var buf bytes.Buffer
	err := dataformats.DescribeAndPrint(d, &buf)
	if err != nil {
		t.Fatalf("DescribeAndPrint failed: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "secret_value_1") || strings.Contains(out, "secret_value_2") {
		t.Errorf("Describe output leaked private values: %s", out)
	}
}
