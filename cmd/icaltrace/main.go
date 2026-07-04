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
	"pimtrace/funcs"

	"github.com/arran4/go-evaluator"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
		fmt.Println("Usage: ", os.Args[0], "[Flags]", "[Query]")
		f.PrintDefaults()
		PrintQueryHelp(os.Stdout, *parser)
	}

	if *versionFlag {
		fmt.Println(version, commit, date)
		return
	}

	if err := f.Parse(os.Args[1:]); err != nil {
		log.Printf("Error parsing flags: %s", err)
		os.Exit(-1)
	}

	if *helpFlag || len(os.Args) <= 1 {
		fmt.Println("No query found")
		f.Usage()
		os.Exit(-1)
	}

	data, err := InputHandler(*inputType, *inputFile, os.Stdout)
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
	if err := OutputHandler(data, *outputType, *outputFile); err != nil {
		log.Printf("Write Error: %s", err)
		os.Exit(-1)
	}
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
		_, _ = fmt.Fprintln(w, "Basic Parser")
		_, _ = fmt.Fprintln(w, "")
		_, _ = fmt.Fprintln(w, "Queries can do the following:")
		_, _ = fmt.Fprintln(w, "- Filtering out data")
		_, _ = fmt.Fprintln(w, "- Selecting components of the data to view")
		_, _ = fmt.Fprintln(w, "- Grouping and summarizing data")
		_, _ = fmt.Fprintln(w, "")
		_, _ = fmt.Fprintln(w, "The Queries can be build up like this:")
		_, _ = fmt.Fprintln(w, "Simple filter query")
		_, _ = fmt.Fprintln(w, "\tfilter not p.SUMMARY icontains .Report")
		_, _ = fmt.Fprintln(w, "Which filters all events without SUMMARY containing 'Report' ")
		_, _ = fmt.Fprintln(w, "Then you can convert the emails to tabular form using:")
		_, _ = fmt.Fprintln(w, "\tinto table p.LOCATION f.year[p.DUE] f.month[p.DUE]")
		_, _ = fmt.Fprintln(w, "Which creates table of locations, and the due year and due month.")
		_, _ = fmt.Fprintln(w, "If you wanted a summary / count of the lines you can use the summary converter:")
		_, _ = fmt.Fprintln(w, "\tinto summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count")
		_, _ = fmt.Fprintln(w, "Which groups the calendar invites based on location, due date year, due date month, then adds a count")
		_, _ = fmt.Fprintln(w, "If you want to sort you can use:")
		_, _ = fmt.Fprintln(w, "\tsort f.year[p.DUE] f.month[p.DUE]")
		_, _ = fmt.Fprintln(w, "These can be used in any combination and repeated for the desired effect:")
		_, _ = fmt.Fprintln(w, "\tfilter p.SUMMARY icontains .Report into summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count filter c.year-DUE eq 2022 sort p.LOCATION")
		_, _ = fmt.Fprintln(w, "And so forth")
		_, _ = fmt.Fprintln(w, "")
		_, _ = fmt.Fprintln(w, "Notes:")
		_, _ = fmt.Fprintln(w, "- Single word string literals begin with a `.`")
		_, _ = fmt.Fprintln(w, "- Properties are referred to with `p.` once in table form it becomes a column referred to by `c.`")
		_, _ = fmt.Fprintln(w, "- I haven't implemented sub components yet")
		_, _ = fmt.Fprintln(w, "- Extension PRs are welcome and intended")
		_, _ = fmt.Fprintln(w, "- All functions are preceded by `f.`")
		_, _ = fmt.Fprintln(w, "- Once converted to table form, it can not be converted back to ical / ics")
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
