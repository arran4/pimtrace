package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIMain_StdoutRegression(t *testing.T) {
	// Build the CLI tool temporarily
	dir := t.TempDir()
	binPath := filepath.Join(dir, "icaltrace")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build icaltrace for integration test: %v", err)
	}

	t.Run("version flag", func(t *testing.T) {
		cmd := exec.Command(binPath, "-version")
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err != nil {
			t.Fatalf("command execution failed: %v, stderr: %s", err, stderrBuf.String())
		}

		output := stdoutBuf.String()
		if !strings.Contains(output, "dev none unknown") {
			t.Errorf("stdout did not contain expected version data, got: %s", output)
		}
		if strings.Contains(output, "No query found") {
			t.Errorf("stdout unexpectedly contained 'No query found', got: %s", output)
		}
	})

	t.Run("ical stream down pipeline", func(t *testing.T) {
		cmd := exec.Command(binPath, "-parser", "basic", "-input", "-", "-input-type", "ical", "-output", "-", "-output-type", "ical")

		icalInput := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//arran4//golang-ical//EN
BEGIN:VEVENT
UID:12345
SUMMARY:Test Event
END:VEVENT
END:VCALENDAR
`
		cmd.Stdin = strings.NewReader(icalInput)
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err != nil {
			t.Fatalf("command execution failed: %v, stderr: %s", err, stderrBuf.String())
		}

		output := stdoutBuf.String()
		if !strings.Contains(output, "SUMMARY:Test Event") {
			t.Errorf("stdout did not contain expected iCal data, got: %s", output)
		}
	})
}
