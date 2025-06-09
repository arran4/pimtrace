package groupdata

import (
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestRowGetSize(t *testing.T) {
	r := &Row{
		Headers: map[string]int{"a": 0},
		Row:     []pimtrace.Value{pimtrace.SimpleStringValue("x")},
		Contents: tabledata.Data{
			&tabledata.Row{},
			&tabledata.Row{},
			&tabledata.Row{},
		},
	}
	v, err := r.Get("sz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i := v.Integer()
	if i == nil || *i != 3 {
		t.Fatalf("expected 3, got %v", v)
	}
}
