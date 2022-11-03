package ast

import (
	"embed"
	_ "embed"
	"github.com/google/go-cmp/cmp"
	"pimtrace"
	"pimtrace/internal/csvdata"
	"testing"
)

var (
	//go:embed "testdata"
	testdata embed.FS
)

func LoadData1(fn string) pimtrace.Data[*csvdata.CSVRow] {
	f, err := testdata.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r, err := csvdata.ReadCSV(f)
	if err != nil {
		panic(err)
	}
	return r
}

func TestCompoundStatement_Execute(t *testing.T) {
	header1 := map[string]int{"address": 3, "currency": 4, "email": 2, "name": 0, "numberrange": 5, "phone": 1}
	tests := []struct {
		name       string
		Statements Operation[*csvdata.CSVRow]
		data       pimtrace.Data[*csvdata.CSVRow]
		want       pimtrace.Data[*csvdata.CSVRow]
		wantErr    bool
	}{
		{
			name: "Simple filter",
			Statements: &FilterStatement[*csvdata.CSVRow]{
				Expression: &Op[*csvdata.CSVRow]{Op: EqualOp, LHS: EntryExpression[*csvdata.CSVRow]("h.numberrange"), RHS: ConstantExpression[*csvdata.CSVRow]("4")},
			},
			data: LoadData1("testdata/data10.csv"),
			want: csvdata.CSVDataType{
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
