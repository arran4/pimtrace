package funcs

import (
	"fmt"
	"io"
	"strings"
)

func PrintFunctionList(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Functions: ")
	for _, f := range Functions[ValueExpression]() {
		for _, af := range f.Arguments() {
			args := make([]string, 0, len(af.Args))
			for _, aff := range af.Args {
				args = append(args, aff.String())
			}
			fn := fmt.Sprintf("f.%s[%s]", f.Name(), strings.Join(args, ","))
			_, _ = fmt.Fprintf(w, "%-40s%40s\n", fn, af.Description)
		}
	}
}
