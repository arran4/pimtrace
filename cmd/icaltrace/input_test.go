package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	old := os.Stdout
	os.Stdout = w
	f()
	if err := w.Close(); err != nil {
		return "", err
	}
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		_ = r.Close()
		return "", err
	}
	if err := r.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestPrintInputHelpContainsIcal(t *testing.T) {
	out, err := captureOutput(PrintInputHelp)
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "ical") {
		t.Errorf("expected help to contain 'ical' but got %q", out)
	}
}
