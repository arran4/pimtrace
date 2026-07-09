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

type mapFSAdapter struct {
	fs fstest.MapFS
}

type nopFile struct {
	io.Reader
}

func (n nopFile) Close() error { return nil }
func (n nopFile) Write(p []byte) (int, error) { return 0, io.ErrClosedPipe }
func (n nopFile) Seek(offset int64, whence int) (int64, error) { return 0, io.EOF }

func (m mapFSAdapter) OpenFile(name string, flag int, perm os.FileMode) (fsys.File, error) {
	f, err := m.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return nopFile{Reader: f}, nil
}

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
	_, _ = w.WriteString(mailContent)
	_ = w.Close()

	data, err := InputHandler("mailfile", "-")
	if err != nil {
		t.Errorf("InputHandler(mailfile, -) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, -) expected 1 msg")
	}
}

func TestInputHandler_File(t *testing.T) {
	oldFS := fsys.DefaultFS
	defer func() { fsys.DefaultFS = oldFS }()

	mailContent := `From: "John" <john@example.com>
To: "Jane" <jane@example.com>
Subject: Test
Date: Thu, 13 Feb 1969 23:32:54 -0330

body
`
	mockFS := mapFSAdapter{
		fs: fstest.MapFS{
			"test.eml": &fstest.MapFile{Data: []byte(mailContent)},
		},
	}
	fsys.DefaultFS = mockFS

	data, err := InputHandler("mailfile", "test.eml")
	if err != nil {
		t.Errorf("InputHandler(mailfile, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, file) expected 1 msg")
	}
}
