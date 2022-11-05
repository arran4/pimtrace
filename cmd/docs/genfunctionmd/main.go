package main

import (
	"fmt"
	"log"
	"os"
	"pimtrace/ast"
	"pimtrace/funcs"
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
	for _, fun := range funcs.Functions[ast.ValueExpression]() {
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
