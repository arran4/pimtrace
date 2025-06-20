package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}
	r.Close()
	return buf.String()
}

func TestPrintInputHelpContainsIcal(t *testing.T) {
	out := captureOutput(PrintInputHelp)
	if !strings.Contains(out, "ical") {
		t.Errorf("expected help to contain 'ical' but got %q", out)
	}
}
