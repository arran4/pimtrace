package main

import (
	"fmt"
	"log"
	"os"
	"pimtrace/ast"
	"pimtrace/funcs"
	"sort"
	"strings"
)

func main() {
	f, err := os.Create("functions.md")
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	fmt.Fprintln(f, "# Functions")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "| Function Def | Description |")
	fmt.Fprintln(f, "| --- | --- |")
	funcs := funcs.Functions[ast.ValueExpression]()
	funNames := make([]string, 0, len(funcs))
	for funName := range funcs {
		funNames = append(funNames, funName)
	}
	sort.Strings(funNames)
	for _, funName := range funNames {
		fun := funcs[funName]
		for _, af := range fun.Arguments() {
			args := make([]string, 0, len(af.Args))
			for _, aff := range af.Args {
				args = append(args, aff.String())
			}
			fn := fmt.Sprintf("f.%s[%s]", fun.Name(), strings.Join(args, ","))
			fmt.Fprintf(f, "| `%s` | %s |\n", fn, af.Description)
		}
	}
}
