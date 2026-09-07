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

	// We build from the root directory but this test runs in cmd/csvtrace, so we use `go build -o ... .`
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build csvtrace for integration test: %v", err)
	}

	// Create a dummy input CSV file
	inputFile := filepath.Join(dir, "input.csv")
	err := os.WriteFile(inputFile, []byte("Name,Age\nAlice,30\nBob,25\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy input file: %v", err)
	}

	// Run the built binary
	// csvtrace -parser basic -input input.csv -input-type csv -output - -output-type csv
	cmd := exec.Command(binPath, "-parser", "basic", "-input", inputFile, "-input-type", "csv", "-output", "-", "-output-type", "csv")
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	// Important: DO NOT wire cmd.Stdin, or test that it gets closed

	err = cmd.Run()
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdoutBuf.String()
	if !strings.Contains(output, "Name,Age") || !strings.Contains(output, "Alice,30") {
		t.Errorf("stdout did not contain expected CSV data, got: %s", output)
	}
}
