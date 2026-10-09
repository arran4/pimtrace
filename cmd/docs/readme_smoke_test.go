package main_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
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
	csvPath := filepath.Join(dir, "csvtrace")
	mailPath := filepath.Join(dir, "mailtrace")
	icalPath := filepath.Join(dir, "icaltrace")

	if err := exec.Command("go", "build", "-o", csvPath, "../../cmd/csvtrace").Run(); err != nil {
		t.Fatalf("failed to build csvtrace")
	}
	if err := exec.Command("go", "build", "-o", mailPath, "../../cmd/mailtrace").Run(); err != nil {
		t.Fatalf("failed to build mailtrace")
	}
	if err := exec.Command("go", "build", "-o", icalPath, "../../cmd/icaltrace").Run(); err != nil {
		t.Fatalf("failed to build icaltrace")
	}

	// 1. CSVTrace
	out := runCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "summary", "c.Category", "calculate", "f.sum[c.Amount]"})
	if !strings.Contains(out, "Food") || !strings.Contains(out, "150") || !strings.Contains(out, "SUM-AMOUNT") {
		t.Errorf("Unexpected output for csvtrace: %s", out)
	}

	// 2. MailTrace
	out = runCmd(t, mailPath, []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "summary", "h.From", "calculate", "f.count", "sort", "f.count", "desc", "limit", "10"})
	if !strings.Contains(out, "test@example.com") || !strings.Contains(out, "COUNT") {
		t.Errorf("Unexpected output for mailtrace: %s", out)
	}

	// 3. ICalTrace
	out = runCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "filter", "p.SUMMARY", "icontains", ".Meeting", "into", "table", "p.DTSTART", "p.SUMMARY"})
	if !strings.Contains(out, "Meeting") || !strings.Contains(out, "20230101T100000Z") {
		t.Errorf("Unexpected output for icaltrace: %s", out)
	}
}
