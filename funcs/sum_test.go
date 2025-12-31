package funcs

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"pimtrace"
	"pimtrace/dataformats/groupdata"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestSum_Run(t *testing.T) {
	s := Sum[ValueExpression]{}

	// Helper to create a groupdata.Row with n items
	createGroupRow := func(values ...int) *groupdata.Row {
		rows := make([]*tabledata.Row, len(values))
		for i, v := range values {
			rows[i] = &tabledata.Row{
				Headers: map[string]int{"val": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(v)},
			}
		}
		return &groupdata.Row{
			Contents: tabledata.Data(rows),
		}
	}

	for _, test := range []struct {
		Name      string
		Input     pimtrace.Entry
		InputArgs []ValueExpression
		Output    pimtrace.Value
		Err       error
	}{
		{
			Name:      "Sum Integer Column",
			Input:     createGroupRow(10, 20, 30),
			InputArgs: []ValueExpression{EntryExpression("c.val")},
			Output:    pimtrace.SimpleIntegerValue(60),
			Err:       nil,
		},
		{
			Name:      "No Args (Returns 1)",
			Input:     createGroupRow(10, 20, 30),
			InputArgs: []ValueExpression{},
			Output:    pimtrace.SimpleIntegerValue(1),
			Err:       nil,
		},
		{
			Name: "Not a groupdata.Row (Returns 1)",
			Input: &tabledata.Row{
				Headers: map[string]int{"val": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(1)},
			},
			InputArgs: []ValueExpression{EntryExpression("c.val")},
			Output:    pimtrace.SimpleIntegerValue(1),
			Err:       nil,
		},
		{
			Name:      "Empty Group (Returns 0)",
			Input:     createGroupRow(),
			InputArgs: []ValueExpression{EntryExpression("c.val")},
			Output:    pimtrace.SimpleIntegerValue(0),
			Err:       nil,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			res, err := s.Run(test.Input, test.InputArgs)
			if (err != nil) != (test.Err != nil) || (err != nil && !errors.Is(err, test.Err)) {
				if test.Err == nil {
					t.Errorf("Got error when wanted none: %s", err)
				} else if err == nil {
					t.Errorf("Didn't get an error when we were expecting one: %s", test.Err)
				} else {
					t.Errorf("Got %s expected: %s", err, test.Err)
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
