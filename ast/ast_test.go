package ast

import (
	"pimtrace"
	"reflect"
	"testing"
)

func TestCompoundStatement_Execute(t *testing.T) {
	tests := []struct {
		name       string
		Statements []Operation
		data       pimtrace.Data
		want       pimtrace.Data
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CompoundStatement{
				Statements: tt.Statements,
			}
			got, err := o.Execute(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
