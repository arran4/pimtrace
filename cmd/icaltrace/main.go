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
	 if _, err := fmt.Fprintln(w, "This tool is for helping you filter, query and summarize ical files in a comprehensible way"); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, "The usage is as follows:"); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, "\ticaltrace -parser basic -input events.ical -input-type ical -output table $QUERY"); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, "In this example it selects the basic parser, reads from events.ical, of the type ical. Outputs a table"); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, "and runs query $QUERY. You are required to specify all of these arguments."); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
	switch parser {
	case "basic":
		 if _, err := fmt.Fprintln(w, "Basic Parser"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Queries can do the following:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Filtering out data"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Selecting components of the data to view"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Grouping and summarizing data"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "The Queries can be build up like this:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Simple filter query"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "\tfilter not p.SUMMARY icontains .Report"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Which filters all events without SUMMARY containing 'Report' "); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Then you can convert the emails to tabular form using:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "\tinto table p.LOCATION f.year[p.DUE] f.month[p.DUE]"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Which creates table of locations, and the due year and due month."); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "If you wanted a summary / count of the lines you can use the summary converter:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "\tinto summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Which groups the calendar invites based on location, due date year, due date month, then adds a count"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "If you want to sort you can use:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "\tsort f.year[p.DUE] f.month[p.DUE]"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "These can be used in any combination and repeated for the desired effect:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "\tfilter p.SUMMARY icontains .Report into summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count filter c.year-DUE eq 2022 sort p.LOCATION"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "And so forth"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "Notes:"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Single word string literals begin with a `.`"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Properties are referred to with `p.` once in table form it becomes a column referred to by `c.`"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- I haven't implemented sub components yet"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Extension PRs are welcome and intended"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- All functions are preceded by `f.`"); err != nil { panic(err) }
		 if _, err := fmt.Fprintln(w, "- Once converted to table form, it can not be converted back to ical / ics"); err != nil { panic(err) }
	}
	 if _, err := fmt.Fprintln(w, "A complete list of functions supported:"); err != nil { panic(err) }
	// TODO funcs.PrintFunctionList(w) when updated
	funcs.PrintFunctionList()
	 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, "List of supported input types:"); err != nil { panic(err) }
	PrintInputHelp(w)
	 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
	 if _, err := fmt.Fprintln(w, "List of supported output types: (Must be supported based on query.)"); err != nil { panic(err) }
	// TODO dataformats.PrintOutputHelp(w, customOutputs) when updated
	dataformats.PrintOutputHelp(customOutputs)
	 if _, err := fmt.Fprintln(w, ""); err != nil { panic(err) }
}
