package main_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runCookbookCmd(t *testing.T, binPath string, args []string) string {
	cmd := exec.Command(binPath, args...)
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

	tests := []struct {
		name     string
		bin      string
		args     []string
		expected []string
	}{
		// CSV Examples
		{
			name:     "CSV Select Columns",
			bin:      csvPath,
			args:     []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "table", "c.Date", "c.Amount"},
			expected: []string{"DATE", "AMOUNT", "2023-01-01"},
		},
		{
			name:     "CSV Filter Rows",
			bin:      csvPath,
			args:     []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "filter", "c.Amount", "gt", ".50", "into", "table", "c.Category", "c.Amount"},
			expected: []string{"CATEGORY", "AMOUNT", "Food", "Utilities", "100", "120"},
		},
		{
			name:     "CSV Summarize",
			bin:      csvPath,
			args:     []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "into", "summary", "c.Category", "calculate", "f.sum[c.Amount]"},
			expected: []string{"SUM-AMOUNT", "150", "120", "45"},
		},
		{
			name:     "CSV Sort",
			bin:      csvPath,
			args:     []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "sort", "c.Amount", "desc", "into", "table", "c.Date", "c.Category", "c.Amount"},
			expected: []string{"AMOUNT"},
		},
		{
			name:     "CSV Range",
			bin:      csvPath,
			args:     []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-parser", "basic", "-output-type", "table", "filter", "c.Amount", "gt", ".10", "and", "c.Amount", "lte", ".100", "into", "table", "c.Date", "c.Category", "c.Amount"},
			expected: []string{"AMOUNT"},
		},

		// Mail Examples
		{
			name:     "Mail Find Subjects Senders",
			bin:      mailPath,
			args:     []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "table", "h.From", "h.Subject"},
			expected: []string{"FROM", "SUBJECT", "sender"},
		},
		{
			name:     "Mail Top Senders",
			bin:      mailPath,
			args:     []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "summary", "h.From", "calculate", "f.count", "sort", "c.count", "desc", "limit", "10"},
			expected: []string{"FROM", "COUNT", "11", "10"},
		},
		{
			name:     "Mail Year Month Summaries",
			bin:      mailPath,
			args:     []string{"-input", "testdata/inbox.mbox", "-input-type", "mbox", "-parser", "basic", "-output-type", "table", "into", "summary", "f.year[h.Date]", "f.month[h.Date]", "calculate", "f.count"},
			expected: []string{"YEAR", "MONTH", "COUNT"},
		},

		// iCal Examples
		{
			name:     "iCal Filter Events",
			bin:      icalPath,
			args:     []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "filter", "p.SUMMARY", "icontains", ".Meeting", "into", "table", "p.DTSTART", "p.SUMMARY"},
			expected: []string{"DTSTART", "SUMMARY", "meEtinG"},
		},
		{
			name:     "iCal Monthly Summaries",
			bin:      icalPath,
			args:     []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "summary", "f.year[p.DTSTART]", "f.month[p.DTSTART]", "calculate", "f.count"},
			expected: []string{"YEAR", "MONTH", "COUNT"},
		},
		{
			name:     "iCal Weekday Hour Duration",
			bin:      icalPath,
			args:     []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "summary", "f.weekday[p.DTSTART]", "f.hour[p.DTSTART]", "calculate", "f.count", "sort", "c.count", "desc", "limit", "5"},
			expected: []string{"WEEKDAY", "HOUR", "COUNT"},
		},
		{
			name:     "iCal Longest Meetings",
			bin:      icalPath,
			args:     []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "sort", "f.duration", "desc", "limit", "5", "into", "table", "p.SUMMARY", "f.duration"},
			expected: []string{"SUMMARY", "DURATION"},
		},
		{
			name:     "iCal Total Meeting Time per Day",
			bin:      icalPath,
			args:     []string{"-input", "testdata/calendar.ics", "-input-type", "ical", "-parser", "basic", "-output-type", "table", "into", "summary", "f.date[p.DTSTART]", "calculate", "f.sum[f.duration]", "sort", "c.date-DTSTART", "asc", "into", "table", "c.date-DTSTART", "c.sum-duration"},
			expected: []string{"DATE-DTSTART", "SUM-DURATION"},
		},

		// Discover Workflow
		{
			name:     "CSV Describe",
			bin:      csvPath,
			args:     []string{"-input", "testdata/expenses.csv", "-input-type", "csv", "-describe"},
			expected: []string{"Expression", "c.Amount", "c.Category"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := runCookbookCmd(t, tt.bin, tt.args)
			for _, exp := range tt.expected {
				if !strings.Contains(out, exp) {
					t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", exp, out)
				}
			}
		})
	}
}
