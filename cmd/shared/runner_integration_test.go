package shared

import (
	"bytes"
	"io"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/maildata"
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
