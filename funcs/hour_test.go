package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestHour_Run(t *testing.T) {
	fn := Hour[ValueExpression]{}

	// Test a time with hour 10
	res, err := fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27T10:30:00Z")},
	}, nil)
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 10 {
		t.Errorf("Run() expected 10, got %v", res)
	}

	// Test invalid
	_, err = fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("invalid")},
	}, nil)
	if err == nil {
		t.Errorf("Run() expected error for invalid input")
	}
}

func TestHour_NameAndArguments(t *testing.T) {
	fn := Hour[ValueExpression]{}
	if fn.Name() != "hour" {
		t.Errorf("Name() = %v, want hour", fn.Name())
	}
	if len(fn.Arguments()) != 2 {
		t.Errorf("Arguments() len = %v, want 2", len(fn.Arguments()))
	}
}
