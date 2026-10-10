package main_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runCookbookCmd(t *testing.T, binPath string, args []string, stdinStr ...string) string {
	cmd := exec.Command(binPath, args...)
	if len(stdinStr) > 0 {
		cmd.Stdin = strings.NewReader(stdinStr[0])
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s\nargs: %v", err, stderr.String(), args)
	}
	return stdout.String()
}

func TestCookbookSmoke(t *testing.T) {
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

	t.Run("CSV Sort Check", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "sort", "c.Amount", "desc", "into", "table", "c.Amount"})
		parsed := parseTable(out)
		if len(parsed) < 4 {
			t.Fatalf("Expected at least 3 data rows in testdata/expenses.csv, got %v", len(parsed))
		}
		if parsed[0][0] != "AMOUNT" {
			t.Fatalf("Expected header AMOUNT, got %v", parsed[0][0])
		}
		if len(parsed) < 5 {
			t.Fatalf("Expected 4 data rows, got %v", len(parsed)-1)
		}
		if parsed[1][0] != "50" || parsed[2][0] != "45" || parsed[3][0] != "120" || parsed[4][0] != "100" {
			t.Errorf("Sorting failed, expected 50, 45, 120, 100, got %v, %v, %v, %v", parsed[1][0], parsed[2][0], parsed[3][0], parsed[4][0])
		}
	})

	t.Run("CSV Range Check", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "filter", "c.Amount", "gt", ".10", "and", "c.Amount", "lte", ".100", "into", "table", "c.Amount"})
		parsed := parseTable(out)
		// Should include 100, 45, 50
		has100 := false
		has120 := false
		for _, row := range parsed {
			if row[0] == "100" {
				has100 = true
			}
			if row[0] == "120" {
				has120 = true
			}
		}
		if len(parsed) != 4 {
			t.Errorf("Expected 3 matched rows + 1 header, got %d", len(parsed))
		}
		if !has100 {
			t.Errorf("Range filter incorrectly excluded 100")
		}
		if has120 {
			t.Errorf("Range filter failed to exclude 120")
		}
	})

	t.Run("CSV Output Mode", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "csv", "into", "summary", "c.Category", "calculate", "f.sum[c.Amount]"})
		if !strings.Contains(strings.ToLower(out), "category,sum-amount") {
			t.Errorf("Expected CSV header Category,sum-Amount, got:\n%s", out)
		}
		if !strings.Contains(out, "Food,150") {
			t.Errorf("Expected CSV to contain Food sum 150, got:\n%s", out)
		}
	})

	t.Run("Pipeline Stdin Check", func(t *testing.T) {
		// Mock cat testdata/expenses.csv | csvtrace -input - ...
		csvContent := "Date,Category,Amount\n2023-10-25,Coffee,4\n2023-10-25,Laptop,1200\n"
		out := runCookbookCmd(t, csvPath, []string{"-input", "-", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "filter", "c.Amount", "gt", ".10", "into", "table", "c.Category", "c.Amount"}, csvContent)
		parsed := parseTable(out)
		if len(parsed) != 2 { // 1 header + 1 row (Laptop)
			t.Fatalf("Expected 2 rows (header + 1 matched row), got %v", len(parsed))
		}
		if parsed[1][0] != "Laptop" || parsed[1][1] != "1200" {
			t.Errorf("Pipeline filter failed, expected Laptop 1200, got %v", parsed[1])
		}
	})

	t.Run("iCal Timezone Aware Check", func(t *testing.T) {
		// Use funcs/testdata/timezone.ics to test timezone conversion explicitly
		out := runCookbookCmd(t, icalPath, []string{"-input", "../../funcs/testdata/timezone.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "table", "f.date[p.DTSTART]"})
		parsed := parseTable(out)
		// Check that the date function successfully extracts a local representation
		if len(parsed) < 3 {
			t.Fatalf("Expected results for timezone.ics")
		}
		if parsed[1][0] != "2023-10-27" || parsed[2][0] != "2023-11-10" {
			t.Errorf("Expected 2023-10-27 and 2023-11-10 respecting NY local TZ, got %v and %v", parsed[1][0], parsed[2][0])
		}
	})

	// Restore CSV Examples
	t.Run("CSV Select Columns", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "table", "c.Date", "c.Amount"})
		parsed := parseTable(out)
		if len(parsed) != 5 { // 1 header + 4 rows
			t.Fatalf("Expected 5 rows, got %v", len(parsed))
		}
		if parsed[0][0] != "DATE" || parsed[0][1] != "AMOUNT" {
			t.Errorf("Expected DATE, AMOUNT header")
		}
	})

	t.Run("CSV Filter Rows", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "filter", "c.Amount", "gt", ".50", "into", "table", "c.Category", "c.Amount"})
		parsed := parseTable(out)
		if len(parsed) != 3 { // 1 header + 100, 120
			t.Fatalf("Expected 3 rows, got %v", len(parsed))
		}
		if parsed[1][0] != "Food" || parsed[1][1] != "100" {
			t.Errorf("Expected Food 100")
		}
	})

	t.Run("CSV Summarize", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "summary", "c.Category", "calculate", "f.sum[c.Amount]"})
		parsed := parseTable(out)
		hasFood150 := false
		for _, row := range parsed {
			if row[0] == "Food" && row[1] == "150" {
				hasFood150 = true
			}
		}
		if !hasFood150 {
			t.Errorf("Expected Food to sum to 150")
		}
	})

	// Restore Mail Examples
	t.Run("Mail Find Subjects Senders", func(t *testing.T) {
		out := runCookbookCmd(t, mailPath, []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "table", "h.From", "h.Subject"})
		parsed := parseTable(out)
		if len(parsed) < 3 {
			t.Fatalf("Expected results")
		}
		if parsed[0][0] != "FROM" || parsed[0][1] != "SUBJECT" {
			t.Errorf("Expected FROM, SUBJECT header")
		}
	})

	t.Run("Mail Top Senders", func(t *testing.T) {
		out := runCookbookCmd(t, mailPath, []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "summary", "h.From", "calculate", "f.count", "sort", "c.count", "desc", "limit", "10"})
		parsed := parseTable(out)
		if len(parsed) != 11 { // 1 header + 10 rows
			t.Fatalf("Expected 11 rows, got %v", len(parsed))
		}
		if parsed[1][0] != "sender11@example.com" || parsed[1][1] != "11" {
			t.Errorf("Expected sender11 to have count 11 at top, got %v", parsed[1])
		}
	})

	t.Run("Mail Year Month Summaries", func(t *testing.T) {
		out := runCookbookCmd(t, mailPath, []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "summary", "f.year[h.Date]", "f.month[h.Date]", "calculate", "f.count"})
		parsed := parseTable(out)
		if len(parsed) < 2 {
			t.Fatalf("Expected results")
		}
		if parsed[0][0] != "YEAR-DATE" || parsed[0][1] != "MONTH-DATE" || parsed[0][2] != "COUNT" {
			t.Errorf("Expected YEAR-DATE, MONTH-DATE, COUNT header, got %v", parsed[0])
		}
	})

	// Restore iCal Examples
	t.Run("iCal Filter Events", func(t *testing.T) {
		out := runCookbookCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "filter", "p.SUMMARY", "icontains", ".Meeting", "into", "table", "p.DTSTART", "p.SUMMARY"})
		parsed := parseTable(out)
		if len(parsed) != 2 {
			t.Fatalf("Expected 2 rows (header + 1), got %v", len(parsed))
		}
		if parsed[1][1] != "meEtinG" {
			t.Errorf("Expected meEtinG, got %v", parsed[1][1])
		}
	})

	t.Run("iCal Monthly Summaries", func(t *testing.T) {
		out := runCookbookCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "summary", "f.year[p.DTSTART]", "f.month[p.DTSTART]", "calculate", "f.count"})
		parsed := parseTable(out)
		if len(parsed) < 2 {
			t.Fatalf("Expected results")
		}
		if parsed[0][0] != "YEAR-DTSTART" {
			t.Errorf("Expected YEAR header, got %v", parsed[0])
		}
	})

	t.Run("iCal Weekday Hour Duration", func(t *testing.T) {
		out := runCookbookCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "summary", "f.weekday[p.DTSTART]", "f.hour[p.DTSTART]", "calculate", "f.count", "sort", "c.count", "desc", "limit", "5"})
		parsed := parseTable(out)
		if len(parsed) < 3 {
			t.Fatalf("Expected results")
		}
		if parsed[0][0] != "WEEKDAY-DTSTART" {
			t.Errorf("Expected WEEKDAY header, got %v", parsed[0])
		}
	})

	t.Run("iCal Longest Meetings", func(t *testing.T) {
		out := runCookbookCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "sort", "f.duration", "desc", "limit", "5", "into", "table", "p.SUMMARY", "f.duration"})
		parsed := parseTable(out)
		if len(parsed) < 3 {
			t.Fatalf("Expected results")
		}
		if parsed[0][0] != "SUMMARY" || parsed[0][1] != "DURATION" {
			t.Errorf("Expected SUMMARY, DURATION header")
		}
	})

	t.Run("iCal Total Meeting Time per Day", func(t *testing.T) {
		out := runCookbookCmd(t, icalPath, []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "summary", "f.date[p.DTSTART]", "calculate", "f.sum[f.duration]", "sort", "c.date-DTSTART", "asc", "into", "table", "c.date-DTSTART", "c.sum-duration"})
		parsed := parseTable(out)
		if len(parsed) < 3 {
			t.Fatalf("Expected results")
		}
		if parsed[0][0] != "DATE-DTSTART" || parsed[0][1] != "SUM-DURATION" {
			t.Errorf("Expected DATE-DTSTART, SUM-DURATION header")
		}
	})

	// Restore Discover Workflow
	t.Run("Discover CSV", func(t *testing.T) {
		out := runCookbookCmd(t, csvPath, []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-describe"})
		if !strings.Contains(out, "c.Amount") || !strings.Contains(out, "c.Category") {
			t.Errorf("Expected describe output to contain c.Amount and c.Category")
		}
	})
}
