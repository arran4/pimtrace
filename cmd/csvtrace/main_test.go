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
	binPath := filepath.Join(dir, "csvtrace")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build csvtrace for integration test: %v", err)
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

	t.Run("csv stream down pipeline", func(t *testing.T) {
		// Run the built binary using stdin input and stdout output
		// echo "Name,Age\nAlice,30\nBob,25\n" | csvtrace -parser basic -input - -input-type csv -output - -output-type csv
		cmd := exec.Command(binPath, "-parser", "basic", "-input", "-", "-input-type", "csv", "-output", "-", "-output-type", "csv")

		cmd.Stdin = strings.NewReader("Name,Age\nAlice,30\nBob,25\n")
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err != nil {
			t.Fatalf("command execution failed: %v, stderr: %s", err, stderrBuf.String())
		}

		output := stdoutBuf.String()
		if !strings.Contains(output, "Name,Age") || !strings.Contains(output, "Alice,30") {
			t.Errorf("stdout did not contain expected CSV data, got: %s", output)
		}
	})

	t.Run("missing input file", func(t *testing.T) {
		cmd := exec.Command(binPath, "-parser", "basic", "-input", "non_existent_file_12345.csv", "-input-type", "csv")
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err == nil {
			t.Fatalf("command execution expected to fail for missing file")
		}

		if cmd.ProcessState.ExitCode() != 1 {
			t.Errorf("expected exit code 1 for runtime read error, got %d", cmd.ProcessState.ExitCode())
		}

		if stdoutBuf.Len() > 0 {
			t.Errorf("stdout expected empty on failure, got: %s", stdoutBuf.String())
		}

		if !strings.Contains(stderrBuf.String(), "Read Error:") {
			t.Errorf("stderr missing 'Read Error:', got: %s", stderrBuf.String())
		}
	})
}

func TestCLIMain_NumericRelationalAndBooleanAcceptance(t *testing.T) {
	dir := t.TempDir()
	binPath := filepath.Join(dir, "csvtrace")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build csvtrace for acceptance tests: %v", err)
	}

	runCSV := func(fixtureContent string, query ...string) (stdout string, stderr string, exitCode int, err error) {
		tmpFile := filepath.Join(t.TempDir(), "data.csv")
		if writeErr := os.WriteFile(tmpFile, []byte(fixtureContent), 0644); writeErr != nil {
			t.Fatalf("failed to write fixture file: %v", writeErr)
		}

		args := append([]string{
			"-parser", "basic",
			"-input", tmpFile,
			"-input-type", "csv",
			"-output", "-",
			"-output-type", "csv",
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

	t.Run("2 < 10 numerically", func(t *testing.T) {
		fixture := "Name,Score\nAlpha,2\nBeta,10\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Score", "lt", ".10")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "Alpha,2") {
			t.Errorf("stdout expected to contain Alpha,2, got: %s", stdout)
		}
		if strings.Contains(stdout, "Beta,10") {
			t.Errorf("stdout unexpectedly contained Beta,10 for < 10, got: %s", stdout)
		}
	})

	t.Run("10 > 2 numerically", func(t *testing.T) {
		fixture := "Name,Score\nAlpha,2\nBeta,10\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Score", "gt", ".2")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "Beta,10") {
			t.Errorf("stdout expected to contain Beta,10, got: %s", stdout)
		}
		if strings.Contains(stdout, "Alpha,2") {
			t.Errorf("stdout unexpectedly contained Alpha,2 for > 2, got: %s", stdout)
		}
	})

	t.Run("decimals", func(t *testing.T) {
		fixture := "Name,Amount\nItemA,2.50\nItemB,10.25\nItemC,100.5\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Amount", "gt", ".5.0", "and", "c.Amount", "lt", ".50.0")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "ItemB,10.25") {
			t.Errorf("stdout expected to contain ItemB,10.25, got: %s", stdout)
		}
		if strings.Contains(stdout, "ItemA,2.50") || strings.Contains(stdout, "ItemC,100.5") {
			t.Errorf("stdout unexpectedly contained excluded items, got: %s", stdout)
		}
	})

	t.Run("negative values", func(t *testing.T) {
		fixture := "Name,Val\nNegSmall,-5\nNegLarge,-20\nPos,5\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Val", "lt", ".-10")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "NegLarge,-20") {
			t.Errorf("stdout expected to contain NegLarge,-20, got: %s", stdout)
		}
		if strings.Contains(stdout, "NegSmall,-5") || strings.Contains(stdout, "Pos,5") {
			t.Errorf("stdout unexpectedly contained excluded items, got: %s", stdout)
		}
	})

	t.Run("zero", func(t *testing.T) {
		fixture := "Name,Val\nZero,0\nNonZero,1\nNeg,-1\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Val", "gte", ".0", "and", "c.Val", "lte", ".0")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "Zero,0") {
			t.Errorf("stdout expected to contain Zero,0, got: %s", stdout)
		}
		if strings.Contains(stdout, "NonZero,1") || strings.Contains(stdout, "Neg,-1") {
			t.Errorf("stdout unexpectedly contained non-zero rows, got: %s", stdout)
		}
	})

	t.Run("valid numeric whitespace", func(t *testing.T) {
		fixture := "Name,Val\nSpaced,\" 5 \"\nOther,10\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Val", "gt", ".4", "and", "c.Val", "lt", ".6")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "Spaced") {
			t.Errorf("stdout expected to contain row with spaced numeric value, got: %s", stdout)
		}
		if strings.Contains(stdout, "Other") {
			t.Errorf("stdout unexpectedly contained Other, got: %s", stdout)
		}
	})

	t.Run("numeric relational expression encountering invalid text errors", func(t *testing.T) {
		fixture := "Name,Score\nAlpha,2\nBad,apple\n"
		_, stderr, code, err := runCSV(fixture, "filter", "c.Score", "gt", ".5")
		if err == nil {
			t.Fatalf("expected command failure on invalid numeric coercion, but succeeded")
		}
		if code != 1 {
			t.Errorf("expected exit code 1 for execution error, got %d", code)
		}
		if !strings.Contains(stderr, "Execute Error:") {
			t.Errorf("stderr expected to contain 'Execute Error:', got: %s", stderr)
		}
		if !strings.Contains(stderr, "numeric comparison failed") && !strings.Contains(stderr, "cannot be coerced") {
			t.Errorf("stderr expected useful coercion error message, got: %s", stderr)
		}
	})

	t.Run("empty or missing numeric value errors", func(t *testing.T) {
		fixture := "Name,Score\nAlpha,2\nMissing,\n"
		_, stderr, code, err := runCSV(fixture, "filter", "c.Score", "gt", ".5")
		if err == nil {
			t.Fatalf("expected command failure on empty numeric value, but succeeded")
		}
		if code != 1 {
			t.Errorf("expected exit code 1 for execution error, got %d", code)
		}
		if !strings.Contains(stderr, "Execute Error:") {
			t.Errorf("stderr expected to contain 'Execute Error:', got: %s", stderr)
		}
	})

	t.Run("legacy textual equality regression 01 vs 1", func(t *testing.T) {
		fixture := "Name,Code\nItemA,01\nItemB,1\n"
		stdout, stderr, code, err := runCSV(fixture, "filter", "c.Code", "eq", ".1")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "ItemB,1") {
			t.Errorf("stdout expected to contain ItemB,1, got: %s", stdout)
		}
		if strings.Contains(stdout, "ItemA,01") {
			t.Errorf("stdout unexpectedly matched 01 for eq .1: equality must remain textual, got: %s", stdout)
		}

		stdout01, stderr01, code01, err01 := runCSV(fixture, "filter", "c.Code", "eq", ".01")
		if err01 != nil {
			t.Fatalf("run failed: %v, stderr: %s", err01, stderr01)
		}
		if code01 != 0 {
			t.Errorf("expected exit code 0, got %d", code01)
		}
		if !strings.Contains(stdout01, "ItemA,01") {
			t.Errorf("stdout expected to contain ItemA,01, got: %s", stdout01)
		}
		if strings.Contains(stdout01, "ItemB,1") {
			t.Errorf("stdout unexpectedly matched 1 for eq .01, got: %s", stdout01)
		}
	})

	t.Run("representative mixed relational and boolean expression", func(t *testing.T) {
		fixture := "Name,Score,Weight\nAlpha,2,10.25\nBeta,10,2.50\nGamma,12,15.00\nDelta,1,1.00\n"
		// ( Score > 5 or Weight > 5.0 ) and not Name eq Beta
		stdout, stderr, code, err := runCSV(fixture, "filter", "(", "c.Score", "gt", ".5", "or", "c.Weight", "gt", ".5.0", ")", "and", "not", "c.Name", "eq", ".Beta")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "Alpha,2,10.25") {
			t.Errorf("stdout expected to contain Alpha, got: %s", stdout)
		}
		if !strings.Contains(stdout, "Gamma,12,15.00") {
			t.Errorf("stdout expected to contain Gamma, got: %s", stdout)
		}
		if strings.Contains(stdout, "Beta,10,2.50") {
			t.Errorf("stdout unexpectedly contained Beta (should be excluded by not Name eq Beta), got: %s", stdout)
		}
		if strings.Contains(stdout, "Delta,1,1.00") {
			t.Errorf("stdout unexpectedly contained Delta (neither score nor weight matched), got: %s", stdout)
		}
	})
}
