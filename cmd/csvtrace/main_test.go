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
