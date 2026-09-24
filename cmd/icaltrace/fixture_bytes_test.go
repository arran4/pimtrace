package main

import (
	"strings"
	"testing"
)

// Preserve the exact line endings exercised by the original CLI acceptance inputs.
func TestEmbeddedAcceptanceFixtureLineEndings(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		crlf    bool
	}{
		{"date range", acceptanceDateRangeICS, true},
		{"all day", acceptanceAllDayICS, true},
		{"duration", acceptanceDurationICS, true},
		{"filter function", acceptanceFilterFunctionICS, true},
		{"simple", acceptanceSimpleICS, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.crlf {
				if !strings.HasSuffix(tc.content, "\r\n") || strings.Contains(strings.ReplaceAll(tc.content, "\r\n", ""), "\n") {
					t.Fatal("expected CRLF on every line, including the final line")
				}
			} else if !strings.HasSuffix(tc.content, "\n") || strings.Contains(tc.content, "\r") {
				t.Fatal("expected LF line endings, including the final line")
			}
		})
	}
}
