package groupdata

import (
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestRowGetSize(t *testing.T) {
	contents := tabledata.Data{
		&tabledata.Row{Headers: map[string]int{"col": 0}, Row: []pimtrace.Value{pimtrace.SimpleStringValue("x")}},
		&tabledata.Row{Headers: map[string]int{"col": 0}, Row: []pimtrace.Value{pimtrace.SimpleStringValue("y")}},
	}
	r := &Row{
		Headers:  map[string]int{"g": 0},
		Row:      []pimtrace.Value{pimtrace.SimpleStringValue("aa")},
		Contents: contents,
	}
	v, err := r.Get("sz")
	if err != nil {
		t.Fatalf("sz err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 2 {
		t.Errorf("sz expected 2 got %v", v)
	}
	v, err = r.Get("sized.c.g")
	if err != nil {
		t.Fatalf("sized err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 2 {
		t.Errorf("sized expected 2 got %v", v)
	}
}
