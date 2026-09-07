package main

import (
	"bytes"
	"os"
	"pimtrace/dataformats/icaldata"
	"strings"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
)

func TestOutputHandlerStreams(t *testing.T) {
	cal := ics.NewCalendar()
	event := cal.AddEvent("12345")
	event.SetSummary("Test Event")
	event.SetStartAt(time.Now())

	data := icaldata.Data{
		&icaldata.ICalWithSource{
			Component:     event,
			ComponentBase: &event.ComponentBase,
			SourceFile:    "test.ics",
		},
	}

	t.Run("ical to stream", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := OutputHandler(data, "ical", "-", buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), "SUMMARY:Test Event") {
			t.Errorf("output does not contain expected ical summary, got: %s", buf.String())
		}
	})

	t.Run("ical to named file", func(t *testing.T) {
		f, err := os.CreateTemp("", "icaltrace-output-*.ics")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		f.Close()

		err = OutputHandler(data, "ical", f.Name())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "SUMMARY:Test Event") {
			t.Errorf("named file output does not contain expected ical summary, got: %s", string(b))
		}
	})
}
