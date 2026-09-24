package main

import (
	"strings"
	"testing"
)

func TestEmbeddedAcceptanceMailFixtureLineEndings(t *testing.T) {
	if !strings.HasSuffix(acceptanceMboxFilterMail, "\r\n") || strings.Contains(strings.ReplaceAll(acceptanceMboxFilterMail, "\r\n", ""), "\n") {
		t.Fatal("mbox acceptance input must retain CRLF on every line, including the final line")
	}
	if !strings.HasSuffix(acceptanceSimpleMail, "\n") || strings.Contains(acceptanceSimpleMail, "\r") {
		t.Fatal("simple email acceptance input must retain LF line endings")
	}
}
