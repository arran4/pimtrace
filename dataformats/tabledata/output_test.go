package tabledata

import (
	"bytes"
	"encoding/csv"
	"os"
	"pimtrace"
	"reflect"
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

	var expected [][]string
	if records[0][0] == "A" {
		expected = [][]string{{"A", "B"}, {"1", "2"}, {"3", "4"}}
	} else {
		expected = [][]string{{"B", "A"}, {"2", "1"}, {"4", "3"}}
	}
	if !reflect.DeepEqual(records, expected) {
		t.Errorf("CSV records mismatch.\nGot: %v\nWant: %v", records, expected)
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
	lines := strings.Split(strings.TrimSpace(out), "\n")

	found12 := false
	found34 := false
	for _, line := range lines {
		if strings.Contains(line, "1") && strings.Contains(line, "2") {
			found12 = true
		}
		if strings.Contains(line, "3") && strings.Contains(line, "4") {
			found34 = true
		}
	}
	if !found12 {
		t.Errorf("Row values 1 and 2 not found on the same line in output:\n%s", out)
	}
	if !found34 {
		t.Errorf("Row values 3 and 4 not found on the same line in output:\n%s", out)
	}
}
