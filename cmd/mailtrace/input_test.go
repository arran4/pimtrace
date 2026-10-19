package main

import (
	"bytes"
	"io"
	"pimtrace/dataformats/maildata"
	"pimtrace/fsys/fsystest"
	"strings"
	"testing"
	"testing/fstest"
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
	mailContent := `From: "John" <john@example.com>
To: "Jane" <jane@example.com>
Subject: Test
Date: Thu, 13 Feb 1969 23:32:54 -0330


body
`
	r := strings.NewReader(mailContent)
	data, err := InputHandler("mailfile", "-", r)
	if err != nil {
		t.Errorf("InputHandler(mailfile, -) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, -) expected 1 msg")
	}
	mb := data.Entry(0).(*maildata.MailWithSource).MailBodies
	if len(mb) < 2 {
		t.Fatalf("expected at least 2 mail bodies (HTML, Text), got %d", len(mb))
	}
	b, err := io.ReadAll(mb[1].Reader())
	if err != nil {
		t.Fatalf("error reading mail body: %v", err)
	}
	if string(b) != "\nbody\n" {
		t.Errorf("InputHandler(mailfile, -) body not parsed properly, got: %q", string(b))
	}
}

func TestInputHandler_File(t *testing.T) {
	mailContent := `From: "John" <john@example.com>
To: "Jane" <jane@example.com>
Subject: Test
Date: Thu, 13 Feb 1969 23:32:54 -0330


body
`
	mockFS := fsystest.MapFSAdapter{
		MapFS: fstest.MapFS{
			"test.eml": &fstest.MapFile{Data: []byte(mailContent)},
		},
	}
	data, err := InputHandler("mailfile", "test.eml", mockFS)
	if err != nil {
		t.Errorf("InputHandler(mailfile, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, file) expected 1 msg")
	}
	mb := data.Entry(0).(*maildata.MailWithSource).MailBodies
	if len(mb) < 2 {
		t.Fatalf("expected at least 2 mail bodies (HTML, Text), got %d", len(mb))
	}
	b, err := io.ReadAll(mb[1].Reader())
	if err != nil {
		t.Fatalf("error reading mail body: %v", err)
	}
	if string(b) != "\nbody\n" {
		t.Errorf("InputHandler(mailfile, file) body not parsed properly, got: %q", string(b))
	}
}
