package main

import (
	"bytes"
	"os"
	"pimtrace/dataformats/maildata"
	"strings"
	"testing"

	"github.com/emersion/go-message/mail"
)

func TestOutputHandler(t *testing.T) {
	h := mail.Header{}
	h.SetSubject("Test Subject")

	data := maildata.Data{
		&maildata.MailWithSource{
			MailHeader: h,
			SourceFile: "test.eml",
		},
	}

	t.Run("mailfile to stream", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := OutputHandler(data, "mailfile", "-", buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), "Subject: Test Subject") {
			t.Errorf("output does not contain expected mailfile subject, got: %s", buf.String())
		}
	})

	t.Run("mbox to stream", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := OutputHandler(data, "mbox", "-", buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), "Subject: Test Subject") {
			t.Errorf("output does not contain expected mbox subject, got: %s", buf.String())
		}
	})

	t.Run("mailfile to named file", func(t *testing.T) {
		f, err := os.CreateTemp("", "mailtrace-output-*.eml")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Remove(f.Name())
		}()
		_ = f.Close()

		err = OutputHandler(data, "mailfile", f.Name())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "Subject: Test Subject") {
			t.Errorf("named file output does not contain expected mailfile subject, got: %s", string(b))
		}
	})
}
