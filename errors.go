package pimtrace

import "fmt"

// PartialResultError indicates that a stream was partially parsed,
// and some unreadable or invalid records were intentionally skipped.
// This allows machines or callers to detect best-effort results
// compared to completely successful processing.
type PartialResultError struct {
	Skipped int
	Message string
}

func (e *PartialResultError) Error() string {
	if e.Skipped == 0 {
		return e.Message
	}
	return fmt.Sprintf("partial result: %s (skipped %d records)", e.Message, e.Skipped)
}
