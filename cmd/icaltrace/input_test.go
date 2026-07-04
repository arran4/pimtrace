package main

import (
	"bytes"
	"io"
	"os"
	"pimtrace"
	"strings"
	"testing"
)

func captureOutput(f func(w io.Writer)) (string, error) {
	var buf bytes.Buffer
	f(&buf)
	return buf.String(), nil
}

func TestPrintInputHelpContainsIcal(t *testing.T) {
	out, err := captureOutput(func(w io.Writer) { PrintInputHelp(w) })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "ical") {
		t.Errorf("expected help to contain 'ical' but got %q", out)
	}
}

func TestInputHandler(t *testing.T) {
	// Test 'list'
	var buf bytes.Buffer
	_, err := InputHandler("list", "", &buf)
	if err != nil {
		t.Errorf("InputHandler(list) error: %v", err)
	}
	if !strings.Contains(buf.String(), "ical") {
		t.Errorf("InputHandler(list) output did not contain expected help")
	}

	// Test unsupported type
	_, err = InputHandler("unknown", "", nil)
	if err == nil {
		t.Errorf("InputHandler(unknown) expected error")
	}

	// For input file, we would need to mock or create a real ical file.
	// The problem is `ReadFile` uses a stream mapping that reads from a file path.
	// But passing an invalid file should give an error. Let's just check the error case.
	_, err = InputHandler("ical", "nonexistent.ics", nil)
	if err == nil {
		t.Errorf("InputHandler(ical, nonexistent.ics) expected error")
	}
}

func TestOutputHandler(t *testing.T) {
	// Test unsupported struct format
	var d pimtrace.Data // Data interface without actual struct implementation
	err := OutputHandler(d, "ical", "-")
	if err == nil {
		t.Errorf("OutputHandler unsupported format expected error")
	}

	// Test passing an invalid format down to dataformats.OutputHandler
	err = OutputHandler(nil, "unknown", "-")
	if err == nil {
		t.Errorf("OutputHandler fallback expected error")
	}
}

func TestPrintQueryHelp(t *testing.T) {
	out, err := captureOutput(func(w io.Writer) { PrintQueryHelp(w, "basic") })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "Basic Parser") {
		t.Errorf("expected help to contain 'Basic Parser' but got %q", out)
	}

	out2, err := captureOutput(func(w io.Writer) { PrintQueryHelp(w, "unknown") })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out2, "A complete list of functions supported") {
		t.Errorf("expected generic help to contain 'A complete list of functions supported' but got %q", out2)
	}
}

// A quick struct to pass the ICalFileOutputCapable check
type dummyICalData struct{}

func (d *dummyICalData) WriteICalFile(fName string) error {
	return nil
}

func (d *dummyICalData) WriteICalStream(f io.Writer, fName string) error {
	return nil
}

func (d *dummyICalData) Len() int { return 0 }
func (d *dummyICalData) Entry(n int) pimtrace.Entry { return nil }
func (d *dummyICalData) Truncate(n int) pimtrace.Data { return nil }
func (d *dummyICalData) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data { return nil }
func (d *dummyICalData) NewSelf() pimtrace.Data { return nil }


func TestOutputHandlerICal(t *testing.T) {
	d := &dummyICalData{}

	err := OutputHandler(d, "ical", "-")
	if err != nil {
		t.Errorf("OutputHandler ical - error: %v", err)
	}

	err = OutputHandler(d, "ical", "test.ics")
	if err != nil {
		t.Errorf("OutputHandler ical file error: %v", err)
	}
}

func TestInputHandler_Stdin(t *testing.T) {
	// Let's test the "-" case which reads from stdin
	// We can replace os.Stdin temporarily or just pass "-" and have it fail on stdin reading / parsing
	// golang-ical parsing an empty stdin or whatever might just return empty or error.

	// Temporarily hijack stdin
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Write a minimal valid ical to our fake stdin so it doesn't block
	_, _ = w.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nEND:VCALENDAR\n")
	_ = w.Close()

	data, err := InputHandler("ical", "-", nil)
	if err != nil {
		t.Errorf("InputHandler(ical, -) error: %v", err)
	}
	if data == nil || data.Len() != 0 {
		// no components in our fake calendar
		t.Errorf("InputHandler(ical, -) expected empty data")
	}
}

func TestInputHandler_File(t *testing.T) {
	// Create a temporary file
	f, err := os.CreateTemp("", "test*.ics")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	_, _ = f.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nEND:VCALENDAR\n")
	_ = f.Close()

	data, err := InputHandler("ical", f.Name(), nil)
	if err != nil {
		t.Errorf("InputHandler(ical, file) error: %v", err)
	}
	if data == nil || data.Len() != 0 {
		t.Errorf("InputHandler(ical, file) expected empty data")
	}
}
