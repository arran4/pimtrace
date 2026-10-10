package main_test

import (
	"bytes"
	"fmt"
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
	parsedMail := parseTable(out)
	// Expect 1 header row + 10 data rows = 11 rows
	if len(parsedMail) != 11 {
		t.Errorf("Expected exactly 11 rows (1 header + 10 data), got %d", len(parsedMail))
	} else {
		if parsedMail[0][1] != "COUNT" {
			t.Errorf("Expected COUNT header, got %s", parsedMail[0][1])
		}
		for i := 1; i <= 10; i++ {
			// Row 1 should be sender11 (11 count), Row 10 should be sender02 (2 count)
			expectedSender := fmt.Sprintf("sender%02d@example.com", 12-i)
			expectedCount := fmt.Sprintf("%d", 12-i)

			if parsedMail[i][0] != expectedSender {
				t.Errorf("Expected sender %s at rank %d, got %s", expectedSender, i, parsedMail[i][0])
			}
			if parsedMail[i][1] != expectedCount {
				t.Errorf("Expected count %s at rank %d, got %s", expectedCount, i, parsedMail[i][1])
			}
		}
		// Verify sender01 is completely excluded
		if strings.Contains(out, "sender01@example.com") {
			t.Errorf("Mailtrace output inappropriately contained sender01: %s", out)
		}
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
