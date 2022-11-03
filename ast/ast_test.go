package ast

import (
	"embed"
	_ "embed"
	"github.com/google/go-cmp/cmp"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"testing"
)

var (
	//go:embed "testdata"
	testdata embed.FS
)

func LoadData1(fn string) pimtrace.Data {
	f, err := testdata.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r, err := tabledata.ReadCSV(f)
	if err != nil {
		panic(err)
	}
	return r
}

func Valueify(ss ...any) (result []pimtrace.Value) {
	for _, s := range ss {
		switch s := s.(type) {
		case string:
			result = append(result, pimtrace.SimpleStringValue(s))
			continue
		case int:
			result = append(result, pimtrace.SimpleIntegerValue(s))
			continue
		default:
			panic("unsupported type")
		}
	}
	return
}

func TestCompoundStatement_Execute(t *testing.T) {
	header1 := map[string]int{"address": 3, "currency": 4, "email": 2, "name": 0, "numberrange": 5, "phone": 1}
	header2 := map[string]int{"Name": 0}
	header3 := map[string]int{"Number": 0, "count": 1, "sum-size": 2}
	tests := []struct {
		name       string
		Statements Operation
		data       pimtrace.Data
		want       pimtrace.Data
		wantErr    bool
	}{
		{
			name: "Simple filter",
			Statements: &FilterStatement{
				Expression: &Op{Op: EqualOp, LHS: EntryExpression("h.numberrange"), RHS: ConstantExpression("4")},
			},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{
					Headers: header1,
					Row: Valueify(
						"Jasper Joseph", "(125) 832-4826", "mauris.vestibulum@protonmail.edu",
						"Ap #783-8034 Nunc Street", "$73.44", "4",
					),
				},
			},
			wantErr: false,
		},
		{
			name: "Simple Table column filter",
			Statements: &TableTransformer{
				Columns: []*ColumnExpression{
					{Name: "Name", Operation: EntryExpression("h.name")},
				},
			},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{Headers: header2, Row: Valueify("Jasper Joseph")},
				{Headers: header2, Row: Valueify("Rogan Hopkins")},
				{Headers: header2, Row: Valueify("Shay Cleveland")},
				{Headers: header2, Row: Valueify("Maite Weaver")},
				{Headers: header2, Row: Valueify("Adria Herring")},
				{Headers: header2, Row: Valueify("Laurel Gonzalez")},
				{Headers: header2, Row: Valueify("Jane Bender")},
				{Headers: header2, Row: Valueify("Melinda Barton")},
				{Headers: header2, Row: Valueify("Colorado Sandoval")},
				{Headers: header2, Row: Valueify("Felix Sutton")},
			},
			wantErr: false,
		},
		{
			name: "Simple Table column filter and sort",
			Statements: &CompoundStatement{Statements: []Operation{
				&TableTransformer{
					Columns: []*ColumnExpression{
						{Name: "Name", Operation: EntryExpression("h.name")},
					},
				},
				&SortTransformer{[]ValueExpression{EntryExpression("c.Name")}},
			}},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{Headers: header2, Row: Valueify("Adria Herring")},
				{Headers: header2, Row: Valueify("Colorado Sandoval")},
				{Headers: header2, Row: Valueify("Felix Sutton")},
				{Headers: header2, Row: Valueify("Jane Bender")},
				{Headers: header2, Row: Valueify("Jasper Joseph")},
				{Headers: header2, Row: Valueify("Laurel Gonzalez")},
				{Headers: header2, Row: Valueify("Maite Weaver")},
				{Headers: header2, Row: Valueify("Melinda Barton")},
				{Headers: header2, Row: Valueify("Rogan Hopkins")},
				{Headers: header2, Row: Valueify("Shay Cleveland")},
			},
			wantErr: false,
		},
		{
			name: "Summary Table with all the functions and group by number",
			Statements: &CompoundStatement{Statements: []Operation{
				&GroupTransformer{
					Columns: []*ColumnExpression{
						{Name: "numberrange", Operation: EntryExpression("h.numberrange")},
					},
				},
				&TableTransformer{
					Columns: []*ColumnExpression{
						{Name: "Number", Operation: EntryExpression("h.numberrange")},
						{Name: "count", Operation: &FunctionExpression{Function: "count"}}, //Args: []ValueExpression{EntryExpression("t.contents")}}},
						{Name: "sum-size", Operation: &FunctionExpression{Function: "sum", Args: []ValueExpression{EntryExpression("h.numberrange")}}},
					},
				},
			}},
			data: LoadData1("testdata/data10.csv"),
			want: tabledata.Data{
				{
					Headers: header3,
					Row:     Valueify("4", 1, 4*1),
				},
				{
					Headers: header3,
					Row:     Valueify("9", 6, 6*9),
				},
				{
					Headers: header3,
					Row:     Valueify("7", 1, 1*7),
				},
				{
					Headers: header3,
					Row:     Valueify("6", 1, 1*6),
				},
				{
					Headers: header3,
					Row:     Valueify("1", 1, 1*1),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.Statements.Execute(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Execute() \n%s", diff)
			}
		})
	}
}
