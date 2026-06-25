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
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	data := Data{}
	if err := data.WriteCSVFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteCSVFile returned error: %v", err)
	}
}

func TestWriteCSVStream(t *testing.T) {
	data := Data{}
	tmpFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	if err := data.WriteCSVStream(tmpFile, "file.csv"); err != nil {
		t.Errorf("WriteCSVStream returned error: %v", err)
	}
}

func TestWriteTableFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	data := Data{}
	if err := data.WriteTableFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteTableFile returned error: %v", err)
	}
}

func TestWriteTableStream(t *testing.T) {
	data := Data{}
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	if err := data.WriteTableStream(tmpFile, "file.txt"); err != nil {
		t.Errorf("WriteTableStream returned error: %v", err)
	}
}

func TestWriteMBoxFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.mbox")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	data := Data{}
	if err := data.WriteMBoxFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteMBoxFile returned error: %v", err)
	}
}

func TestWriteMBoxStream(t *testing.T) {
	data := Data{}
	tmpFile, err := os.CreateTemp("", "test*.mbox")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	if err := data.WriteMBoxStream(tmpFile, "file.mbox"); err != nil {
		t.Errorf("WriteMBoxStream returned error: %v", err)
	}
}

func TestWriteMailFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	data := Data{}
	if err := data.WriteMailFile(tmpFile.Name()); err != nil {
		t.Errorf("WriteMailFile returned error: %v", err)
	}
}

func TestWriteMailStream(t *testing.T) {
	data := Data{}
	tmpFile, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	if err := data.WriteMailStream(tmpFile, "file.eml"); err != nil {
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
	tmpFile, err := os.CreateTemp("", "test*.mbox")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	_ = data.WriteMBoxStream(tmpFile, "file.mbox")
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
	tmpFile, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	_ = data.WriteMailStream(tmpFile, "file.eml")
}

func TestWriteMailStream_Multipart(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "test*.eml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	_ = data.WriteMailStream(tmpFile, "file.eml")
}
