package funcs

import (
	"errors"
	"fmt"
	"github.com/arran4/golang-ical"
	"pimtrace"
	"pimtrace/dataformats/icaldata"

	"github.com/arran4/go-evaluator"
)

type Duration[T ValueExpression] struct{}

var _ Function[ValueExpression] = Duration[ValueExpression]{}

func (c Duration[T]) Name() string {
	return "duration"
}

func (c Duration[T]) Arguments() []ArgumentList {
	return []ArgumentList{
		{
			Description: "Returns the duration in integer seconds of an iCalendar event or task.",
		},
	}
}

func (c Duration[T]) Run(d pimtrace.Entry, args []T, ctx *evaluator.Context) (pimtrace.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("duration: no arguments expected")
	}

	icw, ok := d.(*icaldata.ICalWithSource)
	if !ok {
		return nil, fmt.Errorf("duration: entry is not an iCalendar component")
	}

	var hasDuration bool
	var durationProp string
	for _, p := range icw.ComponentBase.Properties {
		if p.IANAToken == string(ics.ComponentPropertyDuration) {
			hasDuration = true
			durationProp = p.Value
			break
		}
	}
	var hasDtEnd bool
	for _, p := range icw.ComponentBase.Properties {
		if p.IANAToken == string(ics.ComponentPropertyDtEnd) {
			hasDtEnd = true
			break
		}
	}
	var hasDue bool
	for _, p := range icw.ComponentBase.Properties {
		if p.IANAToken == string(ics.ComponentPropertyDue) {
			hasDue = true
			break
		}
	}

	if ve, ok := icw.Component.(*ics.VEvent); ok {
		if hasDuration && hasDtEnd {
			return nil, fmt.Errorf("duration: contradictory properties, both DURATION and DTEND present")
		}

		if hasDuration {
			return nil, fmt.Errorf("duration: golang-ical does not expose a safe public way to interpret a DURATION property: %s", durationProp)
		}

		if hasDtEnd {
			start, err := ve.GetStartAt()
			if err != nil {
				return nil, fmt.Errorf("duration: failed to get DTSTART: %w", err)
			}
			end, err := ve.GetEndAt()
			if err != nil {
				return nil, fmt.Errorf("duration: failed to get DTEND: %w", err)
			}
			return pimtrace.SimpleIntegerValue(int(end.Sub(start).Seconds())), nil
		}
	} else if vt, ok := icw.Component.(*ics.VTodo); ok {
		if hasDuration && hasDue {
			return nil, fmt.Errorf("duration: contradictory properties, both DURATION and DUE present")
		}

		if hasDuration {
			return nil, fmt.Errorf("duration: golang-ical does not expose a safe public way to interpret a DURATION property: %s", durationProp)
		}

		if hasDue {
			start, err := vt.GetStartAt()
			if err != nil {
				return nil, fmt.Errorf("duration: failed to get DTSTART: %w", err)
			}
			due, err := vt.GetDueAt()
			if err != nil {
				return nil, fmt.Errorf("duration: failed to get DUE: %w", err)
			}
			return pimtrace.SimpleIntegerValue(int(due.Sub(start).Seconds())), nil
		}
	}

	return nil, errors.New("duration: missing duration information")
}
