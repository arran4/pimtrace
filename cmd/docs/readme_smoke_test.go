package main_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runCmd(t *testing.T, binPath string, args []string) string {
	cmd := exec.Command(binPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s\nargs: %v", err, stderr.String(), args)
	}
	return stdout.String()
}

func TestReadmeSmoke(t *testing.T) {
	dir := t.TempDir()

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	csvPath := filepath.Join(dir, "csvtrace"+ext)
	mailPath := filepath.Join(dir, "mailtrace"+ext)
	icalPath := filepath.Join(dir, "icaltrace"+ext)

	if err := exec.Command("go", "build", "-o", csvPath, "../../cmd/csvtrace").Run(); err != nil {
		t.Fatalf("failed to build csvtrace: %v", err)
	}
	if err := exec.Command("go", "build", "-o", mailPath, "../../cmd/mailtrace").Run(); err != nil {
		t.Fatalf("failed to build mailtrace: %v", err)
	}
	if err := exec.Command("go", "build", "-o", icalPath, "../../cmd/icaltrace").Run(); err != nil {
		t.Fatalf("failed to build icaltrace: %v", err)
	}

	// 1. CSVTrace
	out := runCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "summary", "c.Category", "calculate", "f.sum[c.Amount]"})
	if !strings.Contains(out, "SUM-AMOUNT") || !strings.Contains(out, "Food      |        150") || !strings.Contains(out, "Transport |         45") || !strings.Contains(out, "Utilities |        120") {
		t.Errorf("Unexpected output for csvtrace (expected exact Food 150, Transport 45, Utilities 120 totals): %s", out)
	}

	// 2. MailTrace
	out = runCmd(t, mailPath, []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "summary", "h.From", "calculate", "f.count", "sort", "c.count", "desc", "limit", "10"})
	if !strings.Contains(out, "sender11@example.com") || !strings.Contains(out, "11") || !strings.Contains(out, "COUNT") {
		t.Errorf("Unexpected output for mailtrace: %s", out)
	}
	if strings.Contains(out, "sender01@example.com") {
		t.Errorf("Mailtrace limit 10 failed, output contained excluded sender01: %s", out)
	}

	// Count number of data rows + header (2 border rows, 1 header row, 1 separator row, 10 data rows, 1 bottom border)
	// Because out ends without trailing newline in output of strings.Count, there are 14 newlines.
	if strings.Count(out, "\n") != 14 {
		t.Errorf("Unexpected number of lines for mailtrace (expected 14): %d\nOutput: %s", strings.Count(out, "\n"), out)
	}

	// Assert descending order: sender11 should appear before sender10
	idx11 := strings.Index(out, "sender11")
	idx10 := strings.Index(out, "sender10")
	if idx11 == -1 || idx10 == -1 || idx11 > idx10 {
		t.Errorf("Mailtrace order failed, expected sender11 before sender10. Output: %s", out)
	}

	// 3. ICalTrace
	out = runCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "filter", "p.SUMMARY", "icontains", ".Meeting", "into", "table", "p.DTSTART", "p.SUMMARY"})
	if !strings.Contains(out, "meEtinG") || !strings.Contains(out, "20230101T100000Z") {
		t.Errorf("Unexpected output for icaltrace: %s", out)
	}
	if strings.Contains(out, "NotIt") {
		t.Errorf("icaltrace filter failed, output contained excluded NotIt: %s", out)
	}
}
