package main

import (
	"bytes"
	"os"
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

func TestCLIMain_ICalDateRangeFilteringAcceptance(t *testing.T) {
	dir := t.TempDir()
	binPath := filepath.Join(dir, "icaltrace")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build icaltrace for acceptance tests: %v", err)
	}

	icsFixture := "BEGIN:VCALENDAR\r\n" +
		"VERSION:2.0\r\n" +
		"PRODID:-//arran4//golang-ical//EN\r\n" +
		"BEGIN:VEVENT\r\n" +
		"UID:event-2019-early\r\n" +
		"SUMMARY:Pre-Holiday Review\r\n" +
		"DTSTART:20191230T100000Z\r\n" +
		"DTEND:20191230T110000Z\r\n" +
		"LOCATION:Room A\r\n" +
		"END:VEVENT\r\n" +
		"BEGIN:VEVENT\r\n" +
		"UID:event-2020-target\r\n" +
		"SUMMARY:New Year Strategic Planning\r\n" +
		"DTSTART:20200102T090000Z\r\n" +
		"DTEND:20200102T170000Z\r\n" +
		"LOCATION:Main Auditorium\r\n" +
		"END:VEVENT\r\n" +
		"BEGIN:VEVENT\r\n" +
		"UID:event-2020-late\r\n" +
		"SUMMARY:Sprint Retrospective\r\n" +
		"DTSTART:20200105T140000Z\r\n" +
		"DTEND:20200105T150000Z\r\n" +
		"LOCATION:Room B\r\n" +
		"END:VEVENT\r\n" +
		"END:VCALENDAR\r\n"

	runICal := func(query ...string) (stdout string, stderr string, exitCode int, err error) {
		tmpFile := filepath.Join(t.TempDir(), "calendar.ics")
		if writeErr := os.WriteFile(tmpFile, []byte(icsFixture), 0644); writeErr != nil {
			t.Fatalf("failed to write ics fixture: %v", writeErr)
		}

		args := append([]string{
			"-parser", "basic",
			"-input", tmpFile,
			"-input-type", "ical",
			"-output", "-",
			"-output-type", "ical",
		}, query...)

		cmd := exec.Command(binPath, args...)
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		cmdErr := cmd.Run()
		code := 0
		if cmd.ProcessState != nil {
			code = cmd.ProcessState.ExitCode()
		}
		return outBuf.String(), errBuf.String(), code, cmdErr
	}

	t.Run("date range selects exactly target event", func(t *testing.T) {
		stdout, stderr, code, err := runICal("filter", "p.DTSTART", "gt", ".2020-01-01", "and", "p.DTSTART", "lt", ".2020-01-03")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}

		// Identifying UID and SUMMARY of target event must be present
		if !strings.Contains(stdout, "UID:event-2020-target") {
			t.Errorf("stdout expected to contain target event UID:event-2020-target, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "SUMMARY:New Year Strategic Planning") {
			t.Errorf("stdout expected to contain target event SUMMARY:New Year Strategic Planning, got:\n%s", stdout)
		}

		// Excluded events must be absent
		if strings.Contains(stdout, "UID:event-2019-early") || strings.Contains(stdout, "SUMMARY:Pre-Holiday Review") {
			t.Errorf("stdout unexpectedly contained early event, got:\n%s", stdout)
		}
		if strings.Contains(stdout, "UID:event-2020-late") || strings.Contains(stdout, "SUMMARY:Sprint Retrospective") {
			t.Errorf("stdout unexpectedly contained late event, got:\n%s", stdout)
		}
	})

	t.Run("inclusive date boundary range", func(t *testing.T) {
		stdout, stderr, code, err := runICal("filter", "p.DTSTART", "gte", ".2020-01-02", "and", "p.DTSTART", "lte", ".2020-01-06")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}

		if !strings.Contains(stdout, "UID:event-2020-target") {
			t.Errorf("stdout expected to contain target event UID:event-2020-target, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "UID:event-2020-late") {
			t.Errorf("stdout expected to contain late event UID:event-2020-late, got:\n%s", stdout)
		}
		if strings.Contains(stdout, "UID:event-2019-early") {
			t.Errorf("stdout unexpectedly contained early event, got:\n%s", stdout)
		}
	})

	t.Run("property prefix alias", func(t *testing.T) {
		stdout, stderr, code, err := runICal("filter", "property.DTSTART", "gt", ".2020-01-01", "and", "property.DTSTART", "lt", ".2020-01-03")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "UID:event-2020-target") {
			t.Errorf("stdout expected to contain target event with property. prefix, got:\n%s", stdout)
		}
	})

	t.Run("invalid date coercion produces error", func(t *testing.T) {
		_, stderr, code, err := runICal("filter", "p.DTSTART", "gt", ".not-a-date")
		if err == nil {
			t.Fatalf("expected command failure on invalid date coercion, but succeeded")
		}
		if code != 1 {
			t.Errorf("expected exit code 1, got %d", code)
		}
		if !strings.Contains(stderr, "Execute Error:") {
			t.Errorf("stderr expected to contain 'Execute Error:', got: %s", stderr)
		}
		if !strings.Contains(stderr, "cannot be coerced") && !strings.Contains(stderr, "date comparison failed") {
			t.Errorf("stderr expected useful date coercion error, got: %s", stderr)
		}
	})
}
