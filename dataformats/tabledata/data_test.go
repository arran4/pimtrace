package tabledata

import (
	"pimtrace"
	"testing"
)

func TestRowGetSize(t *testing.T) {
	r := &Row{
		Headers: map[string]int{"a": 0, "b": 1},
		Row:     []pimtrace.Value{pimtrace.SimpleStringValue("x"), pimtrace.SimpleStringValue("y")},
	}
	v, err := r.Get("sz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i := v.Integer()
	if i == nil || *i != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}
