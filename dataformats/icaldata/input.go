package icaldata

import (
	"fmt"
	"github.com/arran4/golang-ical"
	"io"
)

func ReadICalStream(f io.Reader, fType string, fName string, ops ...any) ([]*ICalWithSource, error) {
	var result []*ICalWithSource
	cs, err := ics.ParseCalendar(f)
	if err != nil {
		return nil, fmt.Errorf("ical stream: %w", err)
	}
	for _, ic := range cs.Components {
		var cb *ics.ComponentBase
		switch c := ic.(type) {
		case *ics.VEvent:
			cb = &c.ComponentBase
		case *ics.VTodo:
			cb = &c.ComponentBase
		case *ics.VBusy:
			cb = &c.ComponentBase
		case *ics.VJournal:
			cb = &c.ComponentBase
		}
		header := make(map[string]int, len(cb.Properties))
		for i, p := range cb.Properties {
			header[p.IANAToken] = i
		}
		result = append(result, &ICalWithSource{
			Component:     ic,
			ComponentBase: cb,
			Header:        header,
			SourceFile:    fName,
			SourceType:    fType,
		})
	}
	return result, nil
}
