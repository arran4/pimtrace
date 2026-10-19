package tests

import (
	"bytes"
	"errors"
	"io"
	"pimtrace"
	"pimtrace/cmd/shared"
	"pimtrace/dataformats/maildata"
	"strings"
	"testing"
)

func TestCLI_StrictFlagPropagates(t *testing.T) {
	handler := func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
		strictFound := false
		for _, op := range ops {
			if s, ok := op.(string); ok && s == "strict" {
				strictFound = true
			}
		}
		if !strictFound {
			return nil, errors.New("strict flag not found in ops")
		}
		return maildata.Data{}, nil
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cfg := &shared.Config{
		Stdout:       stdout,
		Stderr:       stderr,
		Stdin:        strings.NewReader(""),
		Args:         []string{"-strict", "-input", "-", "-input-type", "mbox", "-output-type", "count", "-parser", "basic", "into", "table", "c.A"},
		Name:         "testcmd",
		InputHandler: handler,
		OutputHandler: func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error {
			return nil
		},
		PrintQueryHelp: func(w io.Writer, parser string) {},
		PrintVersion:   func(w io.Writer) {},
	}

	exitCode := shared.Run(cfg)
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr.String())
	}
}
