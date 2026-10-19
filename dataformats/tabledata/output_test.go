package tabledata

import (
	"bytes"
	"encoding/csv"
	"os"
	"pimtrace"
	"strings"
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

	cr := csv.NewReader(&buf)
	records, err := cr.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV output: %v", err)
	}

	if len(records) != 3 { // 1 header row + 2 data rows
		t.Fatalf("expected 3 rows, got %d", len(records))
	}

	// Order of headers A and B is non-deterministic since it uses map[string]int, but we check if both are in row 0
	if (records[0][0] != "A" && records[0][1] != "A") || (records[0][0] != "B" && records[0][1] != "B") {
		t.Errorf("expected headers A and B, got %v", records[0])
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

	// Validate table rows minimally without strict whitespace matching
	if !strings.Contains(out, "A") || !strings.Contains(out, "B") {
		t.Errorf("missing headers in output: %s", out)
	}
	if !strings.Contains(out, "1") || !strings.Contains(out, "2") {
		t.Errorf("missing row 1 in output: %s", out)
	}
	if !strings.Contains(out, "3") || !strings.Contains(out, "4") {
		t.Errorf("missing row 2 in output: %s", out)
	}
}
