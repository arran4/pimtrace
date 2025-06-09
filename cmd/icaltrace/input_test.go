package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintInputHelp(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	PrintInputHelp()
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "ical") {
		t.Errorf("expected help to mention ical, got %s", string(out))
	}
}
