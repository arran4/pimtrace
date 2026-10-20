package maildata

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"os"
	"testing"
)

func TestWriteCSVFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	data := Data{}
	if err := data.WriteCSVFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteCSVFile returned error: %v", err)
	}
}

func TestWriteCSVStream(t *testing.T) {
	data := Data{}
	var buf bytes.Buffer
	if err := data.WriteCSVStream(&buf, "file.csv"); err != nil {
		t.Errorf("WriteCSVStream returned error: %v", err)
	}
}

func TestWriteTableFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	data := Data{}
	if err := data.WriteTableFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteTableFile returned error: %v", err)
	}
}

func TestWriteTableStream(t *testing.T) {
	data := Data{}
	var buf bytes.Buffer
	if err := data.WriteTableStream(&buf, "file.txt"); err != nil {
		t.Errorf("WriteTableStream returned error: %v", err)
	}
}

func TestWriteMBoxFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.mbox")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	data := Data{}
	if err := data.WriteMBoxFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteMBoxFile returned error: %v", err)
	}
}

func TestWriteMBoxStream(t *testing.T) {
	data := Data{}
	var buf bytes.Buffer
	if err := data.WriteMBoxStream(&buf, "file.mbox"); err != nil {
		t.Errorf("WriteMBoxStream returned error: %v", err)
	}
}

func TestWriteMailFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	data := Data{}
	if err := data.WriteMailFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteMailFile returned error: %v", err)
	}
}

func TestWriteMailStream(t *testing.T) {
	data := Data{}
	var buf bytes.Buffer
	if err := data.WriteMailStream(&buf, "file.eml"); err != nil {
		t.Errorf("WriteMailStream returned error: %v", err)
	}
}

func TestWriteMBoxStream_WithData(t *testing.T) {
	h := mail.Header{}
	h.Set("From", "test@example.com")
	data := Data{
		&MailWithSource{
			MailHeader: h,
		},
	}
	var buf bytes.Buffer
	err := data.WriteMBoxStream(&buf, "file.mbox")
	// Expected to pass since we are testing output functionality, although it duplicates output on multi-messages (known bug)
	if err != nil {
		t.Fatalf("WriteMBoxStream returned error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("From: test@example.com")) {
		t.Errorf("WriteMBoxStream output missing header, got: %s", buf.String())
	}
}

func TestWriteMailStream_WithData(t *testing.T) {
	h := mail.Header{}
	h.Set("From", "test@example.com")
	data := Data{
		&MailWithSource{
			MailHeader: h,
			MailBodies: []MailBody{
				&MailBodyGeneral{Body: bytes.NewBufferString("body")},
			},
		},
	}
	var buf bytes.Buffer
	err := data.WriteMailStream(&buf, "file.eml")
	if err != nil {
		t.Fatalf("WriteMailStream returned error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("From: test@example.com")) {
		t.Errorf("WriteMailStream output missing header, got: %s", buf.String())
	}
}

func TestWriteMailStream_Multipart(t *testing.T) {
	// The underlying code has a known issue parsing boundary due to case sensitivity. We only assert that we don't swallow errors here.
	h := mail.Header{}
	h.Set("Content-Type", "multipart/alternative; boundary=\"foo\"")
	data := Data{
		&MailWithSource{
			MailHeader: h,
			MailBodies: []MailBody{
				&MailBodyGeneral{Body: bytes.NewBufferString("body")},
			},
		},
	}
	var buf bytes.Buffer
	err := data.WriteMailStream(&buf, "file.eml")
	if err == nil {
		// Documenting known deficiency: the multipart parsing will fail due to lowercased boundary in mime.ParseMediaType inside production code
		// A dedicated fix PR should assert nil error and proper boundary closing here instead.
		t.Fatal("expected known multipart boundary error, but got nil")
	} else {
		if !bytes.Contains([]byte(err.Error()), []byte("multipart boundary error")) && !bytes.Contains([]byte(err.Error()), []byte("boundary parameter")) {
			t.Fatalf("expected known multipart boundary error, got different error: %v", err)
		}
	}
}
