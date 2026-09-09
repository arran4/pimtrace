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
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
			return OutputHandler(data, outputType, outputFile, stdout)
		},
	}))
}

func PrintQueryHelp(w io.Writer, parser string) {
	_, _ = fmt.Fprintln(w, "This tool is for helping you filter, query and summarize ical files in a comprehensible way")
	_, _ = fmt.Fprintln(w, "The usage is as follows:")
	_, _ = fmt.Fprintln(w, "\ticaltrace -parser basic -input events.ical -input-type ical -output table $QUERY")
	_, _ = fmt.Fprintln(w, "In this example it selects the basic parser, reads from events.ical, of the type ical. Outputs a table")
	_, _ = fmt.Fprintln(w, "and runs query $QUERY. You are required to specify all of these arguments.")
	_, _ = fmt.Fprintln(w, "")
	switch parser {
	case "basic":
		_ = basic.PrintHelp(w, "ical")
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
