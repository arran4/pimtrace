package icaldata

import (
	ics "github.com/arran4/golang-ical"
	"testing"
)

func TestICalGetSize(t *testing.T) {
	cb := &ics.ComponentBase{
		Properties: []ics.IANAProperty{
			{BaseProperty: ics.BaseProperty{IANAToken: "SUMMARY", Value: "Hello"}},
			{BaseProperty: ics.BaseProperty{IANAToken: "LOCATION", Value: "Office"}},
		},
	}
	ve := &ics.VEvent{ComponentBase: *cb}
	ic := &ICalWithSource{
		Component:     ve,
		ComponentBase: cb,
		Header:        map[string]int{"SUMMARY": 0, "LOCATION": 1},
	}
	v, err := ic.Get("sz")
	if err != nil {
		t.Fatalf("sz err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 2 {
		t.Errorf("sz expected 2 got %v", v)
	}
	v, err = ic.Get("sized.p.SUMMARY")
	if err != nil {
		t.Fatalf("sized err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 5 {
		t.Errorf("sized expected 5 got %v", v)
	}
}
