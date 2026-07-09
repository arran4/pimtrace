package main

import (
	"bytes"
	"io"
	"os"
	"pimtrace/fsys"
	"strings"
	"testing"
	"testing/fstest"
)





func captureOutput(f func(w io.Writer)) (string, error) {
	var buf bytes.Buffer
	f(&buf)
	return buf.String(), nil
}

func TestPrintInputHelpContainsCsv(t *testing.T) {
	out, err := captureOutput(func(w io.Writer) { PrintInputHelp(w) })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "csv") {
		t.Errorf("expected help to contain 'csv' but got %q", out)
	}
}

func TestInputHandler(t *testing.T) {
	// Test 'list'
	var buf bytes.Buffer
	_, err := InputHandler("list", "", &buf)
	if err != nil {
		t.Errorf("InputHandler(list) error: %v", err)
	}
	if !strings.Contains(buf.String(), "csv") {
		t.Errorf("InputHandler(list) did not print expected help")
	}

	// Test unsupported type
	_, err = InputHandler("unknown", "", nil)
	if err == nil {
		t.Errorf("InputHandler(unknown) expected error")
	}

	// File error case
	_, err = InputHandler("csv", "nonexistent.csv", nil)
	if err == nil {
		t.Errorf("InputHandler(csv, nonexistent.csv) expected error")
	}
}

func TestInputHandler_Stdin(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	_, _ = w.WriteString("col1,col2\nval1,val2\n")
	_ = w.Close()

	data, err := InputHandler("csv", "-", nil)
	if err != nil {
		t.Errorf("InputHandler(csv, -) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(csv, -) expected 1 row")
	}
}

func TestInputHandler_File(t *testing.T) {
	oldFS := fsys.DefaultFS
	defer func() { fsys.DefaultFS = oldFS }()

	mockFS := fsys.MapFSAdapter{
		MapFS: fstest.MapFS{
			"test.csv": &fstest.MapFile{Data: []byte("col1,col2\nval1,val2\n")},
		},
	}
	fsys.DefaultFS = mockFS

	data, err := InputHandler("csv", "test.csv", nil)
	if err != nil {
		t.Errorf("InputHandler(csv, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(csv, file) expected 1 row")
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
