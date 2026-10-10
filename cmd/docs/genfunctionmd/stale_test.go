package main_test

import (
	"bytes"
	"fmt"
	"os"
	"pimtrace/ast"
	"pimtrace/funcs"
	"sort"
	"strings"
	"testing"
)

func TestFunctionsMDStaleness(t *testing.T) {
	// Generate expected markdown in memory
	var buf bytes.Buffer

	_, _ = fmt.Fprintln(&buf, "# Functions")
	_, _ = fmt.Fprintln(&buf, "")
	_, _ = fmt.Fprintln(&buf, "| Function Def | Description |")
	_, _ = fmt.Fprintln(&buf, "| --- | --- |")

	functions := funcs.Functions[ast.ValueExpression]()
	funNames := make([]string, 0, len(functions))
	for funName := range functions {
		funNames = append(funNames, funName)
	}
	sort.Strings(funNames)

	for _, funName := range funNames {
		fun := functions[funName]
		for _, af := range fun.Arguments() {
			args := make([]string, 0, len(af.Args))
			for _, aff := range af.Args {
				args = append(args, aff.String())
			}
			fn := fmt.Sprintf("f.%s[%s]", fun.Name(), strings.Join(args, ","))
			if len(args) == 0 {
				fn = fmt.Sprintf("f.%s", fun.Name())
			}
			_, _ = fmt.Fprintf(&buf, "| `%s` | %s |\n", fn, af.Description)
		}
	}

	expected := buf.String()

	// Read existing functions.md
	existingPath := "../../../functions.md"
	existingBytes, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("Failed to read existing %s: %v", existingPath, err)
	}

	existing := string(existingBytes)

	if expected != existing {
		t.Errorf("functions.md is stale. Please run 'go run cmd/docs/genfunctionmd/main.go' to regenerate it.\n\nDiff check failed.")
	}
}
