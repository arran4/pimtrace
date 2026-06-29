package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func(w io.Writer)) (string, error) {
	var buf bytes.Buffer
	f(&buf)
	return buf.String(), nil
}

func TestPrintInputHelpContainsTypes(t *testing.T) {
	out, err := captureOutput(func(w io.Writer) { PrintInputHelp(w) })
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "mbox") {
		t.Errorf("expected help to contain 'mbox' but got %q", out)
	}
	if !strings.Contains(out, "mailfile") {
		t.Errorf("expected help to contain 'mailfile' but got %q", out)
	}
}

func TestInputHandler(t *testing.T) {
	// Test 'list'
	var buf bytes.Buffer
	_, err := InputHandler("list", "", &buf)
	if err != nil {
		t.Errorf("InputHandler(list) error: %v", err)
	}
	if !strings.Contains(buf.String(), "mbox") {
		t.Errorf("InputHandler(list) output did not contain expected help")
	}

	// Test unsupported type
	_, err = InputHandler("unknown", "")
	if err == nil {
		t.Errorf("InputHandler(unknown) expected error")
	}
}

func TestInputHandler_Stdin(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	mailContent := `From: "John" <john@example.com>
To: "Jane" <jane@example.com>
Subject: Test
Date: Thu, 13 Feb 1969 23:32:54 -0330

body
`
	 if _, err := w.WriteString(mailContent); err != nil { panic(err) }
	if err := w.Close(); err != nil { panic(err) }

	data, err := InputHandler("mailfile", "-")
	if err != nil {
		t.Errorf("InputHandler(mailfile, -) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, -) expected 1 msg")
	}
}

func TestInputHandler_File(t *testing.T) {
	f, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { if err := os.Remove(f.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", f.Name(), err) } })

	mailContent := `From: "John" <john@example.com>
To: "Jane" <jane@example.com>
Subject: Test
Date: Thu, 13 Feb 1969 23:32:54 -0330

body
`
	 if _, err := f.WriteString(mailContent); err != nil { panic(err) }
	if err := f.Close(); err != nil { panic(err) }

	data, err := InputHandler("mailfile", f.Name())
	if err != nil {
		t.Errorf("InputHandler(mailfile, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, file) expected 1 msg")
	}
}
