package main

import (
	"bytes"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/tabledata"
	"strings"
	"testing"
)

func TestOutputHandler(t *testing.T) {
	data := tabledata.Data{
		&tabledata.Row{
			Headers: map[string]int{"Header1": 0, "Header2": 1},
			Row:     []pimtrace.Value{pimtrace.SimpleStringValue("Value1"), pimtrace.SimpleStringValue("Value2")},
		},
	}

	t.Run("csv to stream", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := dataformats.OutputHandler(data, "csv", "-", customOutputs, buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), "Header1,Header2") {
			t.Errorf("output does not contain expected csv headers, got: %s", buf.String())
		}
	})

	t.Run("table to stream", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := dataformats.OutputHandler(data, "table", "-", customOutputs, buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), "HEADER1") {
			t.Errorf("output does not contain expected table headers, got: %s", buf.String())
		}
	})

	t.Run("csv to named file", func(t *testing.T) {
		f, err := os.CreateTemp("", "csvtrace-output-*.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Remove(f.Name())
		}()
		_ = f.Close()

		err = dataformats.OutputHandler(data, "csv", f.Name(), customOutputs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "Header1,Header2") {
			t.Errorf("named file output does not contain expected csv headers, got: %s", string(b))
		}
	})
}
