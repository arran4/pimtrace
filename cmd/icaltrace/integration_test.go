package main

import (
	"bytes"
	"io"
	"pimtrace"
	"pimtrace/cmd/shared"
	"pimtrace/dataformats"
	"strings"
	"testing"
)

func TestDescribeICalIntegration(t *testing.T) {
	icsData := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Example Corp//NONSGML My Product//EN\nBEGIN:VEVENT\nUID:19970901T130000Z-123401@example.com\nDTSTAMP:19970901T130000Z\nDTSTART:19970903T163000Z\nDTEND:19970903T190000Z\nSUMMARY:Secret Meeting\nCLASS:PRIVATE\nX-MY-CUSTOM:Hello\nEND:VEVENT\nEND:VCALENDAR\n"
	r := strings.NewReader(icsData)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cfg := &shared.Config{
		Stdout:       &stdout,
		Stderr:       &stderr,
		Stdin:        r,
		Args:         []string{"-describe", "-input-type", "ical", "-input", "-"},
		Name:         "icaltrace",
		InputHandler: InputHandler,
	}

	code := shared.Run(cfg)
	if code != 0 {
		t.Fatalf("Runner failed with code %d: %s", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "p.SUMMARY") || !strings.Contains(out, "p.X-MY-CUSTOM") {
		t.Errorf("Expected p.SUMMARY and p.X-MY-CUSTOM in describe output: %s", out)
	}
	if strings.Contains(out, "Secret Meeting") {
		t.Errorf("Describe output leaked sensitive values: %s", out)
	}

	// Verify executing a query with the discovered expression
	r2 := strings.NewReader(icsData)
	var stdout2 bytes.Buffer
	var stderr2 bytes.Buffer
	cfg2 := &shared.Config{
		Stdout:       &stdout2,
		Stderr:       &stderr2,
		Stdin:        r2,
		Args:         []string{"-parser", "basic", "-input-type", "ical", "-input", "-", "-output-type", "table", "into", "table", "p.X-MY-CUSTOM"},
		Name:         "icaltrace",
		InputHandler: InputHandler,
		OutputHandler: func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error {
			return dataformats.OutputHandler(data, outputType, outputFile, nil, stdout)
		},
	}
	code2 := shared.Run(cfg2)
	if code2 != 0 {
		t.Fatalf("Runner query failed with code %d: %s", code2, stderr2.String())
	}
	out2 := stdout2.String()
	if !strings.Contains(out2, "Hello") {
		t.Errorf("Expected output to contain 'Hello': %s", out2)
	}
}
