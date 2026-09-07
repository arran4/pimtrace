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
	binPath := filepath.Join(dir, "mailtrace")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build mailtrace for integration test: %v", err)
	}

	mailInput := `From: sender@example.com
To: recipient@example.com
Subject: Test Subject
Date: Mon, 01 Jan 2024 12:00:00 +0000

This is a test email.
`
	mboxInput := "From sender@example.com Mon Jan 01 12:00:00 2024\r\n" + mailInput

	t.Run("mailfile stream down pipeline", func(t *testing.T) {
		cmd := exec.Command(binPath, "-parser", "basic", "-input", "-", "-input-type", "mailfile", "-output", "-", "-output-type", "mailfile")

		cmd.Stdin = strings.NewReader(mailInput)
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err != nil {
			t.Fatalf("command execution failed: %v, stderr: %s", err, stderrBuf.String())
		}

		output := stdoutBuf.String()
		if !strings.Contains(output, "Subject: Test Subject") {
			t.Errorf("stdout did not contain expected mail data, got: %s", output)
		}
	})

	t.Run("mbox stream down pipeline", func(t *testing.T) {
		cmd := exec.Command(binPath, "-parser", "basic", "-input", "-", "-input-type", "mbox", "-output", "-", "-output-type", "mbox")

		cmd.Stdin = strings.NewReader(mboxInput)
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		if err != nil {
			t.Fatalf("command execution failed: %v, stderr: %s", err, stderrBuf.String())
		}

		output := stdoutBuf.String()
		if !strings.Contains(output, "Subject: Test Subject") {
			t.Errorf("stdout did not contain expected mbox data, got: %s", output)
		}
	})
}
