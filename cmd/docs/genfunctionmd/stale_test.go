package main

import (
	"os"
	"testing"
)

func TestFunctionsMDStaleness(t *testing.T) {
	// Generate expected markdown in memory using shared logic
	expected := generateMarkdown()

	// Read existing functions.md
	existingPath := "../../../functions.md"
	existingBytes, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("Failed to read existing %s: %v", existingPath, err)
	}

	existing := string(existingBytes)

	if expected != existing {
		t.Errorf("functions.md is stale. Please run 'go run cmd/docs/genfunctionmd/*.go' to regenerate it.\n\nDiff check failed.")
	}
}
