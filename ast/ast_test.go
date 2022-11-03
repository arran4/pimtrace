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

func TestCompoundStatement_Execute(t *testing.T) {
	header1 := map[string]int{"address": 3, "currency": 4, "email": 2, "name": 0, "numberrange": 5, "phone": 1}
	header2 := map[string]int{"Name": 0}
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
					Row: []string{
						"Jasper Joseph", "(125) 832-4826", "mauris.vestibulum@protonmail.edu",
						"Ap #783-8034 Nunc Street", "$73.44", "4",
					},
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
				{Headers: header2, Row: []string{"Jasper Joseph"}},
				{Headers: header2, Row: []string{"Rogan Hopkins"}},
				{Headers: header2, Row: []string{"Shay Cleveland"}},
				{Headers: header2, Row: []string{"Maite Weaver"}},
				{Headers: header2, Row: []string{"Adria Herring"}},
				{Headers: header2, Row: []string{"Laurel Gonzalez"}},
				{Headers: header2, Row: []string{"Jane Bender"}},
				{Headers: header2, Row: []string{"Melinda Barton"}},
				{Headers: header2, Row: []string{"Colorado Sandoval"}},
				{Headers: header2, Row: []string{"Felix Sutton"}},
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
