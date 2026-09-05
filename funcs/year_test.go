package funcs

import (
	"errors"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/araddon/dateparse"
	"github.com/arran4/go-evaluator"
	"github.com/google/go-cmp/cmp"
)

func TestParseStrict(t *testing.T) {
	for _, test := range []struct {
		Name   string
		Input  string
		Output string
		Err    error
	}{
		{
			Name:   "Should work",
			Input:  "Wed, 8 Feb 2023 19:00:46 +1100 (AEDT)",
			Output: "2023-02-08 08:00:46 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "Different form",
			Input:  "Wed,   8 Feb 2023 19:00:46 +1100",
			Output: "2023-02-08 08:00:46 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "Some sort of date error I got in a google takeout",
			Input:  "Wed,  8 Feb 2023 19:00:46 +1100 (AEDT)",
			Output: "2023-02-08 08:00:46 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "Month out of range",
			Input:  "FRI, 16 AUG 2013  9:39:51 +1000",
			Output: "2013-08-15 23:39:51 +0000 UTC",
			Err:    nil,
		},
		{
			Name:   "GMT-07:00",
			Input:  "Mon, 1 Dec 2008 14:48:22 GMT-07:00",
			Output: "2008-12-01 21:48:22 +0000 UTC",
			Err:    nil,
		},
		//{
		//	Name:   "Replacement character",
		//	Input:  "Sat, 29 Jan 2011 13:54:02 \\xef\\xbf\\xbd+1000",
		//	Output: "2011-01-19 13:39:51 +0000 UTC",
		//	Err:    nil,
		//},
	} {
		t.Run(test.Name, func(t *testing.T) {
			d, err := dateparse.ParseStrict(test.Input)
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
			if diff := cmp.Diff(test.Output, d.UTC().String()); diff != "" {
				t.Errorf("Outputs differ: %s", diff)
			}
		})
	}
}

type EntryExpression string

func (ve EntryExpression) ColumnName() string {
	ss := strings.SplitN(string(ve), ".", 2)
	s := ""
	if len(ss) > 1 {
		s = ss[1]
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return '-'
	}, s)
}

func (ve EntryExpression) Execute(d pimtrace.Entry, ctx *evaluator.Context) (pimtrace.Value, error) {
	return d.Get(string(ve))
}

func TestArg1OnlyToTime(t *testing.T) {
	for _, test := range []struct {
		Name      string
		Input     pimtrace.Entry
		InputArgs []ValueExpression
		Output    string
		Err       error
	}{
		{
			Name: "Replacement character",
			Input: &tabledata.Row{
				Headers: map[string]int{
					"Date": 0,
				},
				Row: []pimtrace.Value{
					pimtrace.SimpleStringValue("Sat, 29 Jan 2011 13:54:02 \xef\xbf\xbd+1000"),
				},
			},
			InputArgs: []ValueExpression{
				EntryExpression("c.Date"),
			},
			Output: "2011-01-29 13:54:02 +0000 UTC",
			Err:    nil,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			d, err := Arg1OnlyToTime("test", test.Input, test.InputArgs, nil)
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
			if diff := cmp.Diff(test.Output, d.UTC().String()); diff != "" {
				t.Errorf("Outputs differ: %s", diff)
			}
		})
	}
}

func TestAsAdapter_Call(t *testing.T) {
	aa := &AsAdapter{}

	// Test normal
	res, err := aa.Call(42, "new_name")
	if err != nil {
		t.Errorf("Call() error = %v", err)
	}
	if res != 42 {
		t.Errorf("Call() expected 42, got %v", res)
	}

	// Test missing args
	_, err = aa.Call()
	if err == nil {
		t.Errorf("Call(empty args) expected error")
	}

	res, err = aa.Call(42)
	if err != nil {
		t.Errorf("Call(1 arg) error = %v", err)
	}
	if res != 42 {
		t.Errorf("Call(1 arg) expected 42, got %v", res)
	}
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
	if err != nil {
		t.Errorf("Run() expected nil error but got %v", err) // Returns nil value instead of bubbling error for Year
	}
	if _, ok := res.(*pimtrace.SimpleNilValue); !ok {
		t.Errorf("Run() empty args expected SimpleNilValue, got %v", res)
	}
}

func TestPrintFunctionList(t *testing.T) {
	PrintFunctionList()
}

func TestArgumentList_String(t *testing.T) {
	al := ArgumentList{
		Args: []Argument{String, Integer, Any},
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

	// Test int64
	res, err = ya.Call(int64(1672531200)) // Jan 1 2023
	if err != nil {
		t.Errorf("Call(int64) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2023 {
		t.Errorf("Call(int64) expected 2023, got %v", res)
	}

	// Test string
	res, err = ya.Call("2024-05-10")
	if err != nil {
		t.Errorf("Call(string) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2024 {
		t.Errorf("Call(string) expected 2024, got %v", res)
	}

	// Test empty string
	res, err = ya.Call("")
	if err != nil {
		t.Errorf("Call(empty string) error = %v", err)
	}
	if res != nil {
		t.Errorf("Call(empty string) expected nil, got %v", res)
	}

	// Test string with symbol
	res, err = ya.Call("2024-05-10+")
	if err != nil {
		t.Errorf("Call(string with symbol) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2024 {
		t.Errorf("Call(string with symbol) expected 2024, got %v", res)
	}

	// Test pimtrace.Value (Integer)
	res, err = ya.Call(pimtrace.SimpleIntegerValue(1672531200))
	if err != nil {
		t.Errorf("Call(pimtrace.Value int) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2023 {
		t.Errorf("Call(pimtrace.Value int) expected 2023, got %v", res)
	}

	// Test pimtrace.Value (String)
	res, err = ya.Call(pimtrace.SimpleStringValue("2023-10-27"))
	if err != nil {
		t.Errorf("Call(pimtrace.Value string) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 2023 {
		t.Errorf("Call(pimtrace.Value string) expected 2023, got %v", res)
	}

	// Test empty args
	_, err = ya.Call()
	if err == nil {
		t.Errorf("Call(empty args) expected error")
	}

	// Test nil arg
	res, err = ya.Call(nil)
	if err != nil {
		t.Errorf("Call(nil arg) error = %v", err)
	}
	if res != nil {
		t.Errorf("Call(nil arg) expected nil, got %v", res)
	}

	// Test unsupported type
	_, err = ya.Call(1.23)
	if err == nil {
		t.Errorf("Call(unsupported type) expected error")
	}
}

func TestMonthAdapter_Call(t *testing.T) {
	ma := &MonthAdapter{}

	// Test int
	res, err := ma.Call(1672531200) // Jan 1 2023
	if err != nil {
		t.Errorf("Call(int) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 1 {
		t.Errorf("Call(int) expected 1, got %v", res)
	}

	// Test int64
	res, err = ma.Call(int64(1672531200)) // Jan 1 2023
	if err != nil {
		t.Errorf("Call(int64) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 1 {
		t.Errorf("Call(int64) expected 1, got %v", res)
	}

	// Test string
	res, err = ma.Call("2024-05-10")
	if err != nil {
		t.Errorf("Call(string) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 5 {
		t.Errorf("Call(string) expected 5, got %v", res)
	}

	// Test empty string
	res, err = ma.Call("")
	if err != nil {
		t.Errorf("Call(empty string) error = %v", err)
	}
	if res != nil {
		t.Errorf("Call(empty string) expected nil, got %v", res)
	}

	// Test string with symbol
	res, err = ma.Call("2024-05-10+")
	if err != nil {
		t.Errorf("Call(string with symbol) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 5 {
		t.Errorf("Call(string with symbol) expected 5, got %v", res)
	}

	// Test pimtrace.Value (Integer)
	res, err = ma.Call(pimtrace.SimpleIntegerValue(1672531200))
	if err != nil {
		t.Errorf("Call(pimtrace.Value int) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 1 {
		t.Errorf("Call(pimtrace.Value int) expected 1, got %v", res)
	}

	// Test pimtrace.Value (String)
	res, err = ma.Call(pimtrace.SimpleStringValue("2023-10-27"))
	if err != nil {
		t.Errorf("Call(pimtrace.Value string) error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 10 {
		t.Errorf("Call(pimtrace.Value string) expected 10, got %v", res)
	}

	// Test empty args
	_, err = ma.Call()
	if err == nil {
		t.Errorf("Call(empty args) expected error")
	}

	// Test nil arg
	res, err = ma.Call(nil)
	if err != nil {
		t.Errorf("Call(nil arg) error = %v", err)
	}
	if res != nil {
		t.Errorf("Call(nil arg) expected nil, got %v", res)
	}

	// Test unsupported type
	_, err = ma.Call(1.23)
	if err == nil {
		t.Errorf("Call(unsupported type) expected error")
	}
}

func TestAs_NameAndArguments(t *testing.T) {
	a := As[ValueExpression]{}
	if n := a.Name(); n != "as" {
		t.Errorf("As.Name() = %v, want as", n)
	}
	args := a.Arguments()
	if len(args) != 1 {
		t.Errorf("As.Arguments() returned %d arguments, want 1", len(args))
	}
}

func TestAs_ColumnName(t *testing.T) {
	a := As[ValueExpression]{}

	// Test normal case
	name := a.ColumnName([]ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("some_value")},
		mockValueExpression{val: pimtrace.SimpleStringValue("new_name")},
	})
	if name != "new_name" {
		t.Errorf("ColumnName() = %v, want new_name", name)
	}

	// Test missing args
	nameEmpty := a.ColumnName([]ValueExpression{
		mockValueExpression{val: pimtrace.SimpleStringValue("some_value")},
	})
	if nameEmpty != "" {
		t.Errorf("ColumnName(missing args) = %v, want empty", nameEmpty)
	}
}

func TestAs_Run(t *testing.T) {
	a := As[ValueExpression]{}
	d := &tabledata.Row{}

	// Test normal
	res, err := a.Run(d, []ValueExpression{
		mockValueExpression{val: pimtrace.SimpleIntegerValue(42)},
		mockValueExpression{val: pimtrace.SimpleStringValue("new_name")},
	}, nil)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if v, ok := res.(pimtrace.SimpleIntegerValue); !ok || int(v) != 42 {
		t.Errorf("Run() expected 42, got %v", res)
	}

	// Test empty args error
	_, err = a.Run(d, []ValueExpression{}, nil)
	if err == nil {
		t.Errorf("Run() empty args expected error, got nil")
	}
}
