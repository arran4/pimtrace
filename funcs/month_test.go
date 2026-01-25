package funcs

import (
	"errors"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMonth_Run(t *testing.T) {
	m := Month[ValueExpression]{}

	for _, test := range []struct {
		Name      string
		Input     pimtrace.Entry
		InputArgs []ValueExpression
		Output    pimtrace.Value
		Err       error
	}{
		{
			Name: "Valid Date String",
			Input: &tabledata.Row{
				Headers: map[string]int{"Date": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleStringValue("2023-02-08")},
			},
			InputArgs: []ValueExpression{EntryExpression("c.Date")},
			Output:    pimtrace.SimpleIntegerValue(2),
			Err:       nil,
		},
		{
			Name: "Valid Unix Timestamp",
			Input: &tabledata.Row{
				Headers: map[string]int{"Date": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(1675843200)}, // 2023-02-08
			},
			InputArgs: []ValueExpression{EntryExpression("c.Date")},
			Output:    pimtrace.SimpleIntegerValue(2),
			Err:       nil,
		},
		{
			Name: "Empty Input",
			Input: &tabledata.Row{
				Headers: map[string]int{"Date": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleStringValue("")},
			},
			InputArgs: []ValueExpression{EntryExpression("c.Date")},
			Output:    &pimtrace.SimpleNilValue{},
			Err:       nil,
		},
		{
			Name: "Nil Input",
			Input: &tabledata.Row{
				Headers: map[string]int{"Date": 0},
				Row:     []pimtrace.Value{&pimtrace.SimpleNilValue{}},
			},
			InputArgs: []ValueExpression{EntryExpression("c.Date")},
			Output:    &pimtrace.SimpleNilValue{},
			Err:       nil,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			res, err := m.Run(test.Input, test.InputArgs, nil)
			if (err != nil) != (test.Err != nil) || (err != nil && !errors.Is(err, test.Err)) {
				// Special case for wrapping error
				if test.Err != nil && err != nil && errors.Is(err, test.Err) {
					// ok
				} else {
					if test.Err == nil {
						t.Errorf("Got error when wanted none: %s", err)
					} else if err == nil {
						t.Errorf("Didn't get an error when we were expecting one: %s", test.Err)
					} else {
						t.Errorf("Got %s expected: %s", err, test.Err)
					}
				}
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(test.Output, res); diff != "" {
				t.Errorf("Outputs differ: %s", diff)
			}
		})
	}
}
