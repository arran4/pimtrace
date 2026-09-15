package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestWeekday_Run(t *testing.T) {
	fn := Weekday[ValueExpression]{}

	// 2023-10-27 is a Friday
	res, err := fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27")},
	}, nil)
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if v, ok := res.(pimtrace.SimpleStringValue); !ok || string(v) != "Friday" {
		t.Errorf("Run() expected Friday, got %v", res)
	}

	// Test invalid
	_, err = fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("invalid")},
	}, nil)
	if err == nil {
		t.Errorf("Run() expected error for invalid input")
	}
}

func TestWeekday_NameAndArguments(t *testing.T) {
	fn := Weekday[ValueExpression]{}
	if fn.Name() != "weekday" {
		t.Errorf("Name() = %v, want weekday", fn.Name())
	}
	if len(fn.Arguments()) != 2 {
		t.Errorf("Arguments() len = %v, want 2", len(fn.Arguments()))
	}
}
