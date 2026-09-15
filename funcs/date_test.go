package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestDate_Run(t *testing.T) {
	fn := Date[ValueExpression]{}

	// Timestamps normalisation on the same local date (different times)
	res1, err1 := fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27T10:30:00Z")},
	}, nil)
	if err1 != nil {
		t.Errorf("Run() unexpected error 1: %v", err1)
	}

	res2, err2 := fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27T23:59:59Z")},
	}, nil)
	if err2 != nil {
		t.Errorf("Run() unexpected error 2: %v", err2)
	}

	if res1.String() != res2.String() || res1.String() != "2023-10-27" {
		t.Errorf("Run() timestamps on same date didn't normalise correctly. Got %v and %v", res1, res2)
	}

	// Test invalid
	_, err := fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("invalid")},
	}, nil)
	if err == nil {
		t.Errorf("Run() expected error for invalid input")
	}
}

func TestDate_NameAndArguments(t *testing.T) {
	fn := Date[ValueExpression]{}
	if fn.Name() != "date" {
		t.Errorf("Name() = %v, want date", fn.Name())
	}
	if len(fn.Arguments()) != 2 {
		t.Errorf("Arguments() len = %v, want 2", len(fn.Arguments()))
	}
}
