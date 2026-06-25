package main

import (
	"bytes"
	"io"
	"os"
	"pimtrace"
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

func TestPrintInputHelpContainsTypes(t *testing.T) {
	out, err := captureOutput(PrintInputHelp)
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	if !strings.Contains(out, "mailfile") || !strings.Contains(out, "mbox") {
		t.Errorf("expected help to contain types but got %q", out)
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

	// Test file missing (should error)
	_, err = InputHandler("mailfile", "nonexistent.eml")
	if err == nil {
		t.Errorf("InputHandler(mailfile, nonexistent.eml) expected error")
	}
	_, err = InputHandler("mbox", "nonexistent.mbox")
	if err == nil {
		t.Errorf("InputHandler(mbox, nonexistent.mbox) expected error")
	}
}

func TestInputHandler_Stdin(t *testing.T) {
	// mailfile stdin
	func() {
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()
		w.WriteString("Subject: Test\n\nBody")
		w.Close()
		_, err := InputHandler("mailfile", "-")
		if err != nil {
			t.Errorf("InputHandler(mailfile, -) error: %v", err)
		}
	}()

	// mbox stdin
	func() {
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()
		w.WriteString("From MAILER-DAEMON Thu Feb 13 23:32:54 1969\nSubject: Test\n\nBody")
		w.Close()
		_, err := InputHandler("mbox", "-")
		if err != nil {
			t.Errorf("InputHandler(mbox, -) error: %v", err)
		}
	}()

	// mboxtar stdin (just testing it parses switch properly and errors out cleanly since it's not a valid tar)
	func() {
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()
		w.WriteString("not a tar")
		w.Close()
		_, err := InputHandler("mboxtar", "-")
		if err == nil {
			t.Errorf("InputHandler(mboxtar, -) expected error on invalid tar")
		}
	}()
}

// A quick struct to pass the MailFileOutputCapable / MBoxOutputCapable checks
type dummyMailData struct{}

func (d *dummyMailData) WriteMailFile(fName string) error { return nil }
func (d *dummyMailData) WriteMailStream(f io.Writer, fName string) error { return nil }
func (d *dummyMailData) WriteMBoxFile(fName string) error { return nil }
func (d *dummyMailData) WriteMBoxStream(f io.Writer, fName string) error { return nil }
func (d *dummyMailData) Len() int { return 0 }
func (d *dummyMailData) Entry(n int) pimtrace.Entry { return nil }
func (d *dummyMailData) Truncate(n int) pimtrace.Data { return nil }
func (d *dummyMailData) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data { return nil }
func (d *dummyMailData) NewSelf() pimtrace.Data { return nil }

func TestOutputHandlerMail(t *testing.T) {
	d := &dummyMailData{}

	err := OutputHandler(d, "mailfile", "-")
	if err != nil {
		t.Errorf("OutputHandler mailfile - error: %v", err)
	}

	err = OutputHandler(d, "mailfile", "test.eml")
	if err != nil {
		t.Errorf("OutputHandler mailfile file error: %v", err)
	}

	err = OutputHandler(d, "mbox", "-")
	if err != nil {
		t.Errorf("OutputHandler mbox - error: %v", err)
	}

	err = OutputHandler(d, "mbox", "test.mbox")
	if err != nil {
		t.Errorf("OutputHandler mbox file error: %v", err)
	}

	// Test unsupported
	var badD pimtrace.Data
	err = OutputHandler(badD, "mailfile", "-")
	if err == nil {
		t.Errorf("OutputHandler expected error for unsupported mailfile")
	}

	err = OutputHandler(badD, "mbox", "-")
	if err == nil {
		t.Errorf("OutputHandler expected error for unsupported mbox")
	}

	err = OutputHandler(nil, "unknown", "-")
	if err == nil {
		t.Errorf("OutputHandler fallback expected error")
	}
}

func TestMain_Help(t *testing.T) {
	// PrintQueryHelp
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

func TestInputHandler_File(t *testing.T) {
	// Create a temporary mail file
	fMail, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(fMail.Name())

	fMail.WriteString("Subject: Test\n\nBody")
	fMail.Close()

	data, err := InputHandler("mailfile", fMail.Name())
	if err != nil {
		t.Errorf("InputHandler(mailfile, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mailfile, file) expected 1 mail")
	}

	// Create a temporary mbox file
	fMbox, err := os.CreateTemp("", "test*.mbox")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(fMbox.Name())

	fMbox.WriteString("From MAILER-DAEMON Thu Feb 13 23:32:54 1969\nSubject: Test\n\nBody\n")
	fMbox.Close()

	data, err = InputHandler("mbox", fMbox.Name())
	if err != nil {
		t.Errorf("InputHandler(mbox, file) error: %v", err)
	}
	if data == nil || data.Len() != 1 {
		t.Errorf("InputHandler(mbox, file) expected 1 mail")
	}

	// mboxtar file (just hitting it for coverage)
	_, _ = InputHandler("mboxtar", fMbox.Name())
}
