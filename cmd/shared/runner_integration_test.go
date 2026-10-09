package shared

import (
	"bytes"
	"fmt"
	"io"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/icaldata"
	"pimtrace/dataformats/maildata"
	"pimtrace/dataformats/tabledata"
	"strings"
	"testing"
)

func TestDescribeEndToEnd(t *testing.T) {
	emlData := "From: alice@example.com\nTo: bob@example.com\nSubject: Secret Subject 123\nX-Custom: hello\n\nBody"
	r := strings.NewReader(emlData)

	msgs, _ := maildata.ReadMailStream(r, "mailfile", "-")
	d := maildata.Data(msgs)
	descs := dataformats.Describe(d, 100)

	var expr string
	for _, desc := range descs {
		if desc.Name == "Subject" {
			expr = desc.Expression
			break
		}
	}

	if expr == "" {
		t.Fatalf("Failed to find Subject expression")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cfg := &Config{
		Stdout: &stdout,
		Stderr: &stderr,
		Args:   []string{"-parser", "basic", "-input-type", "mailfile", "into", "table", expr},
		Name:   "testcmd",
		InputHandler: func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
			return d, nil
		},
		OutputHandler: func(data pimtrace.Data, outputType string, outputFile string, w io.Writer) error {
			if data.Len() != 1 {
				t.Fatalf("Expected 1 row, got %d", data.Len())
			}
			val, _ := data.Entry(0).Get(expr)
			if val == nil || val.String() != "Secret Subject 123" {
				t.Errorf("Expected Subject to evaluate to Secret Subject 123, got %v", val)
			}
			return nil
		},
	}

	code := Run(cfg)
	if code != 0 {
		t.Fatalf("Runner failed with code %d: %s", code, stderr.String())
	}
}

func TestDescribeStdinIntegration(t *testing.T) {
	csvData := "A,B\n1,secret_value\n"
	r := strings.NewReader(csvData)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cfg := &Config{
		Stdout: &stdout,
		Stderr: &stderr,
		Stdin:  r,
		Args:   []string{"-describe", "-input-type", "csv", "-input", "-"},
		Name:   "testcmd",
		InputHandler: func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
			// Actually process the injected stdin mimicking the real command
			var rr io.Reader
			for _, op := range ops {
				if in, ok := op.(struct{ io.Reader }); ok {
					rr = in.Reader
				}
			}
			if rr == nil {
				return nil, fmt.Errorf("stdin missing")
			}
			rows, err := tabledata.ReadCSV(rr, "csv", "-")
			return tabledata.Data(rows), err
		},
	}

	code := Run(cfg)
	if code != 0 {
		t.Fatalf("Runner failed with code %d: %s", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "c.A") || !strings.Contains(out, "c.B") {
		t.Errorf("Expected c.A and c.B in describe output: %s", out)
	}

	if strings.Contains(out, "secret_value") {
		t.Errorf("Describe output leaked sensitive values: %s", out)
	}
}

func TestDescribeICalIntegration(t *testing.T) {
	icsData := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Example Corp//NONSGML My Product//EN\nBEGIN:VEVENT\nUID:19970901T130000Z-123401@example.com\nDTSTAMP:19970901T130000Z\nDTSTART:19970903T163000Z\nDTEND:19970903T190000Z\nSUMMARY:Secret Meeting\nCLASS:PRIVATE\nX-MY-CUSTOM:Hello\nEND:VEVENT\nEND:VCALENDAR\n"
	r := strings.NewReader(icsData)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cfg := &Config{
		Stdout: &stdout,
		Stderr: &stderr,
		Stdin:  r,
		Args:   []string{"-describe", "-input-type", "ical", "-input", "-"},
		Name:   "testcmd",
		InputHandler: func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
			var rr io.Reader
			for _, op := range ops {
				if in, ok := op.(struct{ io.Reader }); ok {
					rr = in.Reader
				}
			}
			comps, err := icaldata.ReadICalStream(rr, "ical", "-")
			return icaldata.Data(comps), err
		},
	}

	code := Run(cfg)
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
}
