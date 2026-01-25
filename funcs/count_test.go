package funcs

import (
	"errors"
	"pimtrace"
	"pimtrace/dataformats/groupdata"
	"pimtrace/dataformats/tabledata"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCount_Run(t *testing.T) {
	c := Count[ValueExpression]{}

	// Helper to create a groupdata.Row with n items
	createGroupRow := func(n int) *groupdata.Row {
		rows := make([]*tabledata.Row, n)
		for i := 0; i < n; i++ {
			rows[i] = &tabledata.Row{
				Headers: map[string]int{"val": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(i)},
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
			Name:      "Count All (No Args)",
			Input:     createGroupRow(5),
			InputArgs: []ValueExpression{},
			Output:    pimtrace.SimpleIntegerValue(5),
			Err:       nil,
		},
		{
			Name:      "Count Truthy (Arg provided)",
			Input:     createGroupRow(5),                           // Values 0, 1, 2, 3, 4
			InputArgs: []ValueExpression{EntryExpression("c.val")}, // Only > 0 is truthy? Integer truthy check: != 0?
			// Checking SimpleIntegerValue.Truthy implementation:
			// func (v SimpleIntegerValue) Truthy() bool { return int(v) != 0 }
			// So 0 is false, 1,2,3,4 are true. Expected count: 4.
			Output: pimtrace.SimpleIntegerValue(4),
			Err:    nil,
		},
		{
			Name: "Not a groupdata.Row",
			Input: &tabledata.Row{
				Headers: map[string]int{"val": 0},
				Row:     []pimtrace.Value{pimtrace.SimpleIntegerValue(1)},
			},
			InputArgs: []ValueExpression{},
			Output:    pimtrace.SimpleIntegerValue(1),
			Err:       nil,
		},
		{
			Name:      "Empty Group",
			Input:     createGroupRow(0),
			InputArgs: []ValueExpression{},
			Output:    pimtrace.SimpleIntegerValue(0),
			Err:       nil,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			res, err := c.Run(test.Input, test.InputArgs, nil)
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
