package tabledata

import (
	"os"
	"pimtrace"
	"testing"
)

func TestWriteCSVFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", tmpFile.Name(), err) } })

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
	t.Cleanup(func() { if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", tmpFile.Name(), err) } })

	if err := data.WriteCSVStream(tmpFile, "file.csv"); err != nil {
		t.Errorf("WriteCSVStream returned error: %v", err)
	}
}

func TestWriteTableFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", tmpFile.Name(), err) } })

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
	t.Cleanup(func() { if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", tmpFile.Name(), err) } })

	if err := data.WriteTableStream(tmpFile, "file.txt"); err != nil {
		t.Errorf("WriteTableStream returned error: %v", err)
	}
}

func TestWriteCSVStream_WithData(t *testing.T) {
	headers := map[string]int{"A": 0, "B": 1}
	data := Data{
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("1"), pimtrace.SimpleStringValue("2")}},
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("3"), pimtrace.SimpleStringValue("4")}},
	}
	tmpFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", tmpFile.Name(), err) } })

	if err := data.WriteCSVStream(tmpFile, "file.csv"); err != nil {
		t.Errorf("WriteCSVStream returned error: %v", err)
	}
}

func TestWriteTableStream_WithData(t *testing.T) {
	headers := map[string]int{"A": 0, "B": 1}
	data := Data{
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("1"), pimtrace.SimpleStringValue("2")}},
		&Row{Headers: headers, Row: []pimtrace.Value{pimtrace.SimpleStringValue("3"), pimtrace.SimpleStringValue("4")}},
	}
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) { t.Logf("failed to remove %v: %v", tmpFile.Name(), err) } })

	if err := data.WriteTableStream(tmpFile, "file.txt"); err != nil {
		t.Errorf("WriteTableStream returned error: %v", err)
	}
}
