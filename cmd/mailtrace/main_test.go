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

func TestCLIMain_MailFilteringBooleanAcceptance(t *testing.T) {
	dir := t.TempDir()
	binPath := filepath.Join(dir, "mailtrace")

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build mailtrace for acceptance tests: %v", err)
	}

	mboxFixture := "From alice@example.com Mon Jan 01 12:00:00 2024\r\n" +
		"From: alice@example.com\r\n" +
		"To: bob@example.com\r\n" +
		"Subject: Urgent Q1 Financial Report\r\n" +
		"Date: Mon, 01 Jan 2024 12:00:00 +0000\r\n" +
		"User-Agent: KMail/5.1\r\n" +
		"\r\n" +
		"Body of Alice message.\r\n" +
		"\r\n" +
		"From charlie@example.com Tue Jan 02 12:00:00 2024\r\n" +
		"From: charlie@example.com\r\n" +
		"To: bob@example.com\r\n" +
		"Subject: Monthly General Newsletter\r\n" +
		"Date: Tue, 02 Jan 2024 12:00:00 +0000\r\n" +
		"User-Agent: Thunderbird/1.0\r\n" +
		"\r\n" +
		"Body of Charlie message.\r\n" +
		"\r\n" +
		"From david@example.com Wed Jan 03 12:00:00 2024\r\n" +
		"From: david@example.com\r\n" +
		"To: alice@example.com\r\n" +
		"Subject: Urgent Security Notice\r\n" +
		"Date: Wed, 03 Jan 2024 12:00:00 +0000\r\n" +
		"User-Agent: Mutt/2.0\r\n" +
		"\r\n" +
		"Body of David message.\r\n" +
		"\r\n" +
		"From eve@example.com Thu Jan 04 12:00:00 2024\r\n" +
		"From: eve@example.com\r\n" +
		"To: charlie@example.com\r\n" +
		"Subject: Team Lunch Plans\r\n" +
		"Date: Thu, 04 Jan 2024 12:00:00 +0000\r\n" +
		"User-Agent: KMail/5.2\r\n" +
		"\r\n" +
		"Body of Eve message.\r\n"

	runMail := func(query ...string) (stdout string, stderr string, exitCode int, err error) {
		tmpFile := filepath.Join(t.TempDir(), "inbox.mbox")
		if writeErr := os.WriteFile(tmpFile, []byte(mboxFixture), 0644); writeErr != nil {
			t.Fatalf("failed to write mbox fixture: %v", writeErr)
		}

		args := append([]string{
			"-parser", "basic",
			"-input", tmpFile,
			"-input-type", "mbox",
			"-output", "-",
			"-output-type", "mbox",
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

	t.Run("and operator", func(t *testing.T) {
		stdout, stderr, code, err := runMail("filter", "h.To", "eq", ".bob@example.com", "and", "h.User-Agent", "icontains", ".kmail")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "alice@example.com") {
			t.Errorf("stdout expected to contain Alice's message, got: %s", stdout)
		}
		if strings.Contains(stdout, "charlie@example.com") || strings.Contains(stdout, "eve@example.com") {
			t.Errorf("stdout unexpectedly contained non-matching messages, got: %s", stdout)
		}
	})

	t.Run("or operator", func(t *testing.T) {
		stdout, stderr, code, err := runMail("filter", "h.From", "eq", ".alice@example.com", "or", "h.From", "eq", ".david@example.com")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "alice@example.com") || !strings.Contains(stdout, "david@example.com") {
			t.Errorf("stdout expected to contain Alice and David messages, got: %s", stdout)
		}
		if strings.Contains(stdout, "charlie@example.com") || strings.Contains(stdout, "eve@example.com") {
			t.Errorf("stdout unexpectedly contained Charlie or Eve messages, got: %s", stdout)
		}
	})

	t.Run("A or B and C precedence", func(t *testing.T) {
		// Evaluated as: From==Alice OR (From==Charlie AND User-Agent==Mutt)
		// Charlie uses Thunderbird (not Mutt). So B and C is false.
		// Alice uses KMail (not Mutt). But Alice matches A!
		// Because 'and' binds tighter than 'or', Alice must be included!
		stdout, stderr, code, err := runMail("filter", "h.From", "eq", ".alice@example.com", "or", "h.From", "eq", ".charlie@example.com", "and", "h.User-Agent", "eq", ".Mutt/2.0")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "alice@example.com") {
			t.Errorf("stdout expected to contain Alice due to and-over-or precedence, got: %s", stdout)
		}
		if strings.Contains(stdout, "charlie@example.com") {
			t.Errorf("stdout unexpectedly contained Charlie, got: %s", stdout)
		}
	})

	t.Run("parentheses overriding precedence", func(t *testing.T) {
		// Evaluated as: (From==Alice OR From==Charlie) AND User-Agent==Mutt
		// Neither Alice nor Charlie uses Mutt, so neither should match!
		stdout, stderr, code, err := runMail("filter", "(", "h.From", "eq", ".alice@example.com", "or", "h.From", "eq", ".charlie@example.com", ")", "and", "h.User-Agent", "eq", ".Mutt/2.0")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if strings.Contains(stdout, "alice@example.com") || strings.Contains(stdout, "charlie@example.com") {
			t.Errorf("stdout unexpectedly contained message when parentheses override precedence, got: %s", stdout)
		}
	})

	t.Run("nested not", func(t *testing.T) {
		stdout, stderr, code, err := runMail("filter", "not", "not", "h.From", "eq", ".alice@example.com")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "alice@example.com") {
			t.Errorf("stdout expected to contain Alice for not not, got: %s", stdout)
		}
		if strings.Contains(stdout, "charlie@example.com") {
			t.Errorf("stdout unexpectedly contained Charlie for not not Alice, got: %s", stdout)
		}
	})

	t.Run("icontains composed with boolean expressions", func(t *testing.T) {
		// Subject contains 'urgent' (case-insensitive) AND NOT from David
		stdout, stderr, code, err := runMail("filter", "h.Subject", "icontains", ".urgent", "and", "not", "h.From", "eq", ".david@example.com")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "alice@example.com") {
			t.Errorf("stdout expected to contain Alice (urgent and not David), got: %s", stdout)
		}
		if strings.Contains(stdout, "david@example.com") {
			t.Errorf("stdout unexpectedly contained David (excluded by not), got: %s", stdout)
		}
		if strings.Contains(stdout, "charlie@example.com") {
			t.Errorf("stdout unexpectedly contained Charlie (newsletter, not urgent), got: %s", stdout)
		}
	})

	t.Run("established mail-header casing and path behavior", func(t *testing.T) {
		// Use lowercased dashed path 'h.user-agent' and 'h.subject'
		stdout, stderr, code, err := runMail("filter", "h.subject", "icontains", ".urgent", "and", "h.user-agent", "icontains", ".kmail")
		if err != nil {
			t.Fatalf("run failed: %v, stderr: %s", err, stderr)
		}
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(stdout, "alice@example.com") {
			t.Errorf("stdout expected to contain Alice with lowercase dashed header paths, got: %s", stdout)
		}
		if strings.Contains(stdout, "david@example.com") || strings.Contains(stdout, "charlie@example.com") {
			t.Errorf("stdout unexpectedly contained other messages, got: %s", stdout)
		}
	})
}
