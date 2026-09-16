package funcs

import (
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

func TestDay_Run(t *testing.T) {
	fn := Day[ValueExpression]{}

	// 2023-10-27 is the 27th day
	res, err := fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27")},
	}, nil)
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 27 {
		t.Errorf("Run() expected 27, got %v", res)
	}

	// Test invalid
	_, err = fn.Run(&tabledata.Row{}, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("invalid")},
	}, nil)
	if err == nil {
		t.Errorf("Run() expected error for invalid input")
	}
}

func TestDay_NameAndArguments(t *testing.T) {
	fn := Day[ValueExpression]{}
	if fn.Name() != "day" {
		t.Errorf("Name() = %v, want day", fn.Name())
	}
	if len(fn.Arguments()) != 2 {
		t.Errorf("Arguments() len = %v, want 2", len(fn.Arguments()))
	}
}
