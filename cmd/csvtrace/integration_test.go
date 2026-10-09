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

func TestDescribeCSVIntegration(t *testing.T) {
	csvData := "A,B\n1,secret_value\n"
	r := strings.NewReader(csvData)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cfg := &shared.Config{
		Stdout:       &stdout,
		Stderr:       &stderr,
		Stdin:        r,
		Args:         []string{"-describe", "-input-type", "csv", "-input", "-"},
		Name:         "csvtrace",
		InputHandler: InputHandler,
	}

	code := shared.Run(cfg)
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

	// Verify executing a query with the discovered expression
	r2 := strings.NewReader(csvData)
	var stdout2 bytes.Buffer
	var stderr2 bytes.Buffer
	cfg2 := &shared.Config{
		Stdout:       &stdout2,
		Stderr:       &stderr2,
		Stdin:        r2,
		Args:         []string{"-parser", "basic", "-input-type", "csv", "-input", "-", "-output-type", "table", "into", "table", "c.A"},
		Name:         "csvtrace",
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
	if !strings.Contains(out2, "1") {
		t.Errorf("Expected output to contain '1': %s", out2)
	}
}
