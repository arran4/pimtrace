package funcs

import (
	"os"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
	"time"

	"github.com/arran4/go-evaluator"
)

type EntryExpression string

func (ve EntryExpression) ColumnName(args []ValueExpression) string {
	return string(ve)
}

func (ve EntryExpression) Execute(d pimtrace.Entry, ctx *evaluator.Context) (pimtrace.Value, error) {
	return d.Get(string(ve))
}

type mockValueExpression struct {
	val pimtrace.Value
	err error
}

func (m mockValueExpression) Execute(d pimtrace.Entry, ctx *evaluator.Context) (pimtrace.Value, error) {
	return m.val, m.err
}

func (m mockValueExpression) ColumnName() string {
	return "mock"
}

func (m mockValueExpression) Evaluate(d interface{}, opts ...any) (interface{}, error) {
	return m.val, m.err
}

func TestYear_NameAndArguments(t *testing.T) {
	y := Year[ValueExpression]{}
	if n := y.Name(); n != "year" {
		t.Errorf("Year.Name() = %v, want year", n)
	}
	args := y.Arguments()
	if len(args) != 2 {
		t.Errorf("Year.Arguments() returned %d arguments, want 2", len(args))
	}
}

func TestYear_Run(t *testing.T) {
	y := Year[ValueExpression]{}

	// Test string parsing
	d := &tabledata.Row{}
	res, err := y.Run(d, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("2023-10-27")},
	}, nil)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2023 {
		t.Errorf("Run() string expected 2023, got %v", res)
	}

	// Test integer parsing (unix time)
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	res, err = y.Run(d, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleIntegerValue(ts)},
	}, nil)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2025 {
		t.Errorf("Run() int expected 2025, got %v", res)
	}

	// Test error case (empty)
	res, err = y.Run(d, []ValueExpression{}, nil)
	if err == nil {
		t.Errorf("Run() expected error but got %v", err)
	}
}

func TestPrintFunctionList(t *testing.T) {
	PrintFunctionList(os.Stdout)
}

func TestArgumentList_String(t *testing.T) {
	al := ArgumentList{
		Args:        []Argument{String, Integer, Any},
		Description: "test desc",
	}
	s := ""
	for i, arg := range al.Args {
		if i > 0 {
			s += ","
		}
		s += arg.String()
	}
	s = "[" + s + "]"

	if s != "[String,Integer,Any]" {
		t.Errorf("ArgumentList formatted = %v, want [String,Integer,Any]", s)
	}
}

func TestYearAdapter_Call(t *testing.T) {
	ya := &YearAdapter{}

	// Test int
	res, err := ya.Call(1672531200) // Jan 1 2023
	if err != nil {
		t.Errorf("Call(int) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2023 {
		t.Errorf("Call(int) expected 2023, got %v", res)
	}

	// Test string
	res, err = ya.Call("2024-05-10")
	if err != nil {
		t.Errorf("Call(string) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2024 {
		t.Errorf("Call(string) expected 2024, got %v", res)
	}

	// Test empty string error
	res, err = ya.Call("")
	if err == nil {
		t.Errorf("Call(empty string) expected error")
	}

	// Test pimtrace.Value (Integer)
	res, err = ya.Call(pimtrace.SimpleIntegerValue(1672531200))
	if err != nil {
		t.Errorf("Call(pimtrace.Value int) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2023 {
		t.Errorf("Call(pimtrace.Value int) expected 2023, got %v", res)
	}
}
