package basic

import (
	"bytes"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	types := []string{"csv", "ical", "mail"}
	for _, ty := range types {
		t.Run(ty, func(t *testing.T) {
			var buf bytes.Buffer
			err := PrintHelp(&buf, ty)
			if err != nil {
				t.Fatalf("PrintHelp failed for type %s: %v", ty, err)
			}
			if buf.Len() == 0 {
				t.Errorf("PrintHelp produced empty output for type %s", ty)
			}
		})
	}
}
