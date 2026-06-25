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

func TestPrintInputHelpContainsCsv(t *testing.T) {
	out, err := captureOutput(PrintInputHelp)
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "csv") {
		t.Errorf("expected help to contain 'csv' but got %q", out)
	}
}

func TestInputHandler(t *testing.T) {
	// Test 'list'
	_, err := InputHandler("list", "")
	if err != nil {
		t.Errorf("InputHandler(list) error: %v", err)
	}

	// Test unsupported type
	_, err = InputHandler("unknown", "")
	if err == nil {
		t.Errorf("InputHandler(unknown) expected error")
	}

	// File error case
	_, err = InputHandler("csv", "nonexistent.csv")
	if err == nil {
		t.Errorf("InputHandler(csv, nonexistent.csv) expected error")
	}
}

func TestInputHandler_Stdin(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	w.WriteString("col1,col2\nval1,val2\n")
	w.Close()

	data, err := InputHandler("csv", "-")
	if err != nil {
		t.Errorf("InputHandler(csv, -) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(csv, -) expected 1 row")
	}
}

func TestInputHandler_File(t *testing.T) {
	// Create a temporary file
	f, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	f.WriteString("col1,col2\nval1,val2\n")
	f.Close()

	data, err := InputHandler("csv", f.Name())
	if err != nil {
		t.Errorf("InputHandler(csv, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(csv, file) expected 1 row")
	}
}

func TestPrintQueryHelp(t *testing.T) {
	out, err := captureOutput(func() { PrintQueryHelp("basic") })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "Basic Parser") {
		t.Errorf("expected help to contain 'Basic Parser' but got %q", out)
	}

	out2, err := captureOutput(func() { PrintQueryHelp("unknown") })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out2, "A complete list of functions supported") {
		t.Errorf("expected generic help to contain 'A complete list of functions supported' but got %q", out2)
	}
}
