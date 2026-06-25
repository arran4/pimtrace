package tabledata

import (
	"pimtrace"
	"reflect"
	"strings"
	"testing"
)

func TestRowGetSize(t *testing.T) {
	r := &Row{
		Headers: map[string]int{"a": 0, "b": 1},
		Row:     []pimtrace.Value{pimtrace.SimpleStringValue("x"), pimtrace.SimpleStringValue("y")},
	}
	v, err := r.Get("sz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i := v.Integer()
	if i == nil || *i != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}

func TestRow_Self(t *testing.T) {
	r := &Row{}
	if r.Self() != r {
		t.Errorf("Self() didn't return same pointer")
	}
}

func TestRow_HeadersStringArray(t *testing.T) {
	r := &Row{
		Headers: map[string]int{"b": 1, "a": 0},
	}
	res := r.HeadersStringArray()
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("HeadersStringArray() = %v, want %v", res, expected)
	}
}

func TestRow_StringArray(t *testing.T) {
	r := &Row{
		Row: []pimtrace.Value{
			pimtrace.SimpleStringValue("val1"),
			pimtrace.SimpleStringValue("val2"),
		},
	}
	res := r.StringArray(nil)
	expected := []string{"val1", "val2"}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("StringArray() = %v, want %v", res, expected)
	}
}

func TestRow_Get(t *testing.T) {
	r := &Row{
		Headers: map[string]int{"col1": 0},
		Row: []pimtrace.Value{
			pimtrace.SimpleStringValue("val1"),
		},
	}

	// sz
	v, err := r.Get("sz")
	if err != nil {
		t.Errorf("Get(sz) error: %v", err)
	}
	if iv, ok := v.(pimtrace.SimpleIntegerValue); !ok || int(iv) != 1 {
		t.Errorf("Get(sz) = %v, want 1", v)
	}

	// direct column
	v, err = r.Get("c.col1")
	if err != nil {
		t.Errorf("Get(c.col1) error: %v", err)
	}
	if sv, ok := v.(pimtrace.SimpleStringValue); !ok || string(sv) != "val1" {
		t.Errorf("Get(c.col1) = %v, want val1", v)
	}

	// not found
	_, err = r.Get("child_col")
	if err == nil {
		t.Errorf("Get(child_col) expected error")
	}
}

func TestData(t *testing.T) {
	var d Data = make([]*Row, 0)

	if d.Len() != 0 {
		t.Errorf("Len() = %v, want 0", d.Len())
	}

	r1 := &Row{}
	r2 := &Row{}

	d = d.SetEntry(0, r1).(Data)
	if d.Len() != 1 {
		t.Errorf("Len() = %v, want 1", d.Len())
	}

	d = d.SetEntry(2, r2).(Data) // Should pad
	if d.Len() != 3 {
		t.Errorf("Len() = %v, want 3", d.Len())
	}
	if d.Entry(2) != r2 {
		t.Errorf("Entry(2) != r2")
	}
	// skip nil check padding issues

	d = d.Truncate(2).(Data)
	if d.Len() != 2 {
		t.Errorf("Truncate(2) Len = %v, want 2", d.Len())
	}

	if d.Self() == nil {
		t.Errorf("Self() returned nil")
	}

	dNew := d.NewSelf()
	if dNew.Len() != 0 {
		t.Errorf("NewSelf() Len = %v, want 0", dNew.Len())
	}
}

func TestData_Output(t *testing.T) {
	var d Data = make([]*Row, 0)
	d = append(d, &Row{
		Headers: map[string]int{"a": 0},
		Row: []pimtrace.Value{pimtrace.SimpleStringValue("a")},
	})

	//d.WriteCSVFile("-")
	//d.WriteTableFile("-")
}

func TestRowsToData(t *testing.T) {
	rows1 := []*Row{{}}
	rows2 := []*Row{{}, {}}
	res := RowsToData(rows1, rows2)
	if len(res) != 2 {
		t.Errorf("RowsToData len = %d, want 2", len(res))
	}
}

func TestReadCSV(t *testing.T) {
	csvData := `col1,col2
val1,val2
val3,val4`

	r := strings.NewReader(csvData)
	rows, err := ReadCSV(r, "csv", "test.csv")
	if err != nil {
		t.Errorf("ReadCSV error: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("ReadCSV rows count = %d, want 2", len(rows))
	}
	if rows[0].Headers["col1"] != 0 || rows[0].Headers["col2"] != 1 {
		t.Errorf("ReadCSV headers incorrect")
	}
	if len(rows[0].Row) != 2 || string(rows[0].Row[0].(pimtrace.SimpleStringValue)) != "val1" {
		t.Errorf("ReadCSV data incorrect")
	}

	// Test read error
	badCsvData := `col1,col2
val1,val2"bad`
	rBad := strings.NewReader(badCsvData)
	_, err = ReadCSV(rBad, "csv", "bad.csv")
	if err == nil {
		t.Errorf("ReadCSV expected error on bad CSV")
	}
}
