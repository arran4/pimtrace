package tabledata

import (
	"bytes"
	"os"
	"pimtrace"
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

func TestWriteCSVStream_WithData(t *testing.T) {
	headers := map[string]int{"A": 0, "B": 1}
	data := Data{
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("1"), pimtrace.SimpleStringValue("2")}},
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("3"), pimtrace.SimpleStringValue("4")}},
	}
	var buf bytes.Buffer
	if err := data.WriteCSVStream(&buf, "file.csv"); err != nil {
		t.Fatalf("WriteCSVStream returned error: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Errorf("WriteCSVStream produced empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("A")) || !bytes.Contains(buf.Bytes(), []byte("B")) {
		t.Errorf("WriteCSVStream output missing headers, got: %s", out)
	}
	if !bytes.Contains(buf.Bytes(), []byte("1")) || !bytes.Contains(buf.Bytes(), []byte("4")) {
		t.Errorf("WriteCSVStream output missing data values, got: %s", out)
	}
}

func TestWriteTableStream_WithData(t *testing.T) {
	headers := map[string]int{"A": 0, "B": 1}
	data := Data{
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("1"), pimtrace.SimpleStringValue("2")}},
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("3"), pimtrace.SimpleStringValue("4")}},
	}
	var buf bytes.Buffer
	if err := data.WriteTableStream(&buf, "file.txt"); err != nil {
		t.Fatalf("WriteTableStream returned error: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Errorf("WriteTableStream produced empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("A")) || !bytes.Contains(buf.Bytes(), []byte("B")) {
		t.Errorf("WriteTableStream output missing headers, got: %s", out)
	}
	if !bytes.Contains(buf.Bytes(), []byte("1")) || !bytes.Contains(buf.Bytes(), []byte("4")) {
		t.Errorf("WriteTableStream output missing data values, got: %s", out)
	}
}
