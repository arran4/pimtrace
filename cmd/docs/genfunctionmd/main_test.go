package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReadmeExamples(t *testing.T) {
	tmpDir := t.TempDir()

	err := exec.Command("go", "build", "-o", tmpDir, "../../...").Run()
	if err != nil {
		t.Fatalf("Failed to build: %v", err)
	}

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	csvtrace := filepath.Join(tmpDir, "csvtrace"+ext)
	mailtrace := filepath.Join(tmpDir, "mailtrace"+ext)
	icaltrace := filepath.Join(tmpDir, "icaltrace"+ext)

	tests := []struct {
		name    string
		command *exec.Cmd
	}{
		{
			name: "CSV Trace Table",
			command: exec.Command(csvtrace, "-input", "../../../dataformats/tabledata/testdata/test.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "table", "c.col1"),
		},
		{
			name: "Mail Trace Table",
			command: exec.Command(mailtrace, "-input", "../../../dataformats/maildata/testdata/test.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "table", "h.Subject"),
		},
		{
			name: "Ical Trace Table",
			command: exec.Command(icaltrace, "-input", "../../../dataformats/icaldata/testdata/meeting_output_example.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "table", "p.SUMMARY"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.command.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
			}
		})
	}
}
