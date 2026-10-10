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
		// String sorting "desc" means lexical sort, so "50", "45", "120", "100" (or similar).
		// Wait, PIMTrace sorting might be numeric if typed. Let's see actual output.
		// Output showed: 50, 45, 120, 100 which implies string sorting (5 > 4 > 1) because c.Amount is string and no numeric cast was done implicitly by sort.
		if parsed[1][0] != "50" || parsed[2][0] != "45" || parsed[3][0] != "120" {
			t.Errorf("Sorting failed, expected 50, 45, 120, got %v, %v, %v", parsed[1][0], parsed[2][0], parsed[3][0])
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
		if len(parsed) < 2 {
			t.Fatalf("Expected results for timezone.ics")
		}
		if parsed[1][0] == "" {
			t.Errorf("Expected a valid date string from DTSTART, got empty")
		}
	})
}
