package main

import (
	"fmt"
	"io"
	"os"
	"pimtrace"
	"pimtrace/argparsers/basic"
	"pimtrace/cmd/shared"
	"pimtrace/dataformats"
	"pimtrace/funcs"
)

var (
	version       = "dev"
	commit        = "none"
	date          = "unknown"
	customOutputs = [][2]string{}
)

func main() {
	os.Exit(shared.Run(&shared.Config{
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		Stdin:          os.Stdin,
		Args:           os.Args[1:],
		Name:           os.Args[0],
		PrintQueryHelp: PrintQueryHelp,
		PrintVersion: func(w io.Writer) {
			_, _ = fmt.Fprintln(w, version, commit, date)
		},
		InputHandler: InputHandler,
		OutputHandler: func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error {
			return dataformats.OutputHandler(data, outputType, outputFile, customOutputs, stdout)
		},
	}))
}

func PrintQueryHelp(w io.Writer, parser string) {
	_, _ = fmt.Fprintln(w, "This tool is for helping you filter, query and summarize CSV/data files in a comprehensible way")
	_, _ = fmt.Fprintln(w, "The usage is as follows:")
	_, _ = fmt.Fprintln(w, "\tcsvtrace -parser basic -input jobs.csv -input-type csv -output table $QUERY")
	_, _ = fmt.Fprintln(w, "In this example it selects the basic parser, reads from jobs.csv, of the type csv. Outputs a table")
	_, _ = fmt.Fprintln(w, "and runs query $QUERY. You are required to specify all of these arguments.")
	_, _ = fmt.Fprintln(w, "")
	switch parser {
	case "basic":
		_ = basic.PrintHelp(w, "csv")
	}
	_, _ = fmt.Fprintln(w, "A complete list of functions supported:")
	funcs.PrintFunctionList(w)
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "List of supported input types:")
	PrintInputHelp(w)
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "List of supported output types: (Must be supported based on query.)")
	dataformats.PrintOutputHelp(w, customOutputs)
	_, _ = fmt.Fprintln(w, "")
}
