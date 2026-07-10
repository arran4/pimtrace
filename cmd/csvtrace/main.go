package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"pimtrace/argparsers/basic"
	"pimtrace/ast"
	"pimtrace/dataformats"
	"pimtrace/fsys"
	"pimtrace/funcs"

	"github.com/arran4/go-evaluator"
)

var (
	version       = "dev"
	commit        = "none"
	date          = "unknown"
	customOutputs = [][2]string{}
)

func main() {
	f := flag.FlagSet{}
	var (
		inputType   = f.String("input-type", "list", "The input type")
		inputFile   = f.String("input", "-", "Input file or - for stdin")
		outputType  = f.String("output-type", "list", "The input type")
		outputFile  = f.String("output", "-", "Output file or - for stdin")
		parser      = f.String("parser", "", "Just use `basic`")
		versionFlag = f.Bool("version", false, "Prints the version")
		helpFlag    = f.Bool("help", false, "Prints help")
	)
	f.Usage = func() {
		_, _ = fmt.Println("Usage: ", os.Args[0], "[Flags]", "[Query]")
		f.PrintDefaults()
		PrintQueryHelp(os.Stdout, *parser)
	}

	if *versionFlag {
		_, _ = fmt.Println(version, commit, date)
		return
	}

	if err := f.Parse(os.Args[1:]); err != nil {
		log.Printf("Error parsing flags: %s", err)
		os.Exit(-1)
	}

	if *helpFlag || len(os.Args) <= 1 {
		_, _ = fmt.Println("No query found")
		f.Usage()
		os.Exit(-1)
	}

	data, err := InputHandler(fsys.OSFS{}, *inputType, *inputFile, os.Stdout)
	if err != nil {
		log.Printf("Read Error: %s", err)
		os.Exit(-1)
	}

	var ops ast.Operation
	switch *parser {
	case "basic":
		ops, err = basic.ParseOperations(f.Args())
		if err != nil {
			log.Printf("Parse Error: %s", err)
			os.Exit(-1)
		}
	default:
		log.Printf("Please use -parser=basic parameter, as maybe one day a more advanced parser will be created")
		os.Exit(-1)
	}

	if ops != nil {
		ctx := &evaluator.Context{
			Functions: map[string]evaluator.Function{
				"year":  &funcs.YearAdapter{},
				"month": &funcs.MonthAdapter{},
				"as":    &funcs.AsAdapter{},
			},
		}
		data, err = ops.Execute(data, ctx)
		if err != nil {
			log.Printf("Execute Error: %s", err)
			os.Exit(-1)
		}
	}
	if err := dataformats.OutputHandler(data, *outputType, *outputFile, customOutputs); err != nil {
		log.Printf("Write Error: %s", err)
		os.Exit(-1)
	}
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
	// TODO funcs.PrintFunctionList(w) when updated
	funcs.PrintFunctionList()
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "List of supported input types:")
	PrintInputHelp(w)
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "List of supported output types: (Must be supported based on query.)")
	// TODO dataformats.PrintOutputHelp(w, customOutputs) when updated
	dataformats.PrintOutputHelp(customOutputs)
	_, _ = fmt.Fprintln(w, "")
}
