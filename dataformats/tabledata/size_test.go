package tabledata

import (
	"pimtrace"
	"testing"
)

func TestRowGetSize(t *testing.T) {
	r := &Row{
		Headers: map[string]int{"name": 0, "value": 1},
		Row: []pimtrace.Value{
			pimtrace.SimpleStringValue("abc"),
			pimtrace.SimpleStringValue("xyz"),
		},
	}
	v, err := r.Get("sz")
	if err != nil {
		t.Fatalf("sz err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 2 {
		t.Errorf("sz expected 2 got %v", v)
	}
	v, err = r.Get("sized.c.name")
	if err != nil {
		t.Fatalf("sized err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 3 {
		t.Errorf("sized expected 3 got %v", v)
	}
}
