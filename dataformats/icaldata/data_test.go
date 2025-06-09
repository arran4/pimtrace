package icaldata

import (
	ics "github.com/arran4/golang-ical"
	"testing"
)

func TestICalWithSourceGetSize(t *testing.T) {
	cb := &ics.ComponentBase{
		Properties: []ics.IANAProperty{
			{BaseProperty: ics.BaseProperty{IANAToken: "SUMMARY"}},
			{BaseProperty: ics.BaseProperty{IANAToken: "LOCATION"}},
		},
	}
	ic := &ICalWithSource{ComponentBase: cb, Header: map[string]int{"SUMMARY": 0, "LOCATION": 1}}
	v, err := ic.Get("sz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i := v.Integer()
	if i == nil || *i != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}
