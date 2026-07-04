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
		progress    = f.Bool("progress", false, "Report progress")
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

	var iops []any

	if *progress {
		iops = append(iops, dataformats.NewProgressor())
	}

	data, err := InputHandler(*inputType, *inputFile, iops...)
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
	_, _ = fmt.Fprintln(w, "This tool is for helping you filter, query and summarize mail files in a comprehensible way")
	fmt.Fprintln(w, "The usage is as follows:")
	fmt.Fprintln(w, "\tmailtrace -parser basic -input mail.mbox -input-type mbox -output table $QUERY")
	fmt.Fprintln(w, "In this example it selects the basic parser, reads from mail.mbox, of the type mbox. Outputs a table")
	fmt.Fprintln(w, "and runs query $QUERY. You are required to specify all of these arguments.")
	fmt.Fprintln(w, "")
	switch parser {
	case "basic":
		fmt.Fprintln(w, "Basic Parser")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Queries can do the following:")
		fmt.Fprintln(w, "- Filtering out data")
		fmt.Fprintln(w, "- Selecting components of the data to view")
		fmt.Fprintln(w, "- Grouping and summarizing data")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "The Queries can be build up like this:")
		fmt.Fprintln(w, "Simple filter query")
		fmt.Fprintln(w, "\tfilter not h.user-agent icontains .Kmail")
		fmt.Fprintln(w, "Which filters all emails sent with a user agent that does not contain `kmail` ")
		fmt.Fprintln(w, "Then you can convert the emails to tabular form using:")
		fmt.Fprintln(w, "\tinto table h.user-agent h.subject f.year[h.date] f.month[h.date]")
		fmt.Fprintln(w, "Which creates table of user agent, subject, and the year and month.")
		fmt.Fprintln(w, "If you wanted a summary / count of the lines you can use the summary converter:")
		fmt.Fprintln(w, "\tinto summary h.user-agent h.subject f.year[h.date] f.month[h.date] calculate f.sum[c.size] f.count")
		fmt.Fprintln(w, "Which groups the mail based on user-agent, subject, year, month, then creates a sum and a count")
		fmt.Fprintln(w, "If you want to sort you can use:")
		fmt.Fprintln(w, "\tsort f.year[h.date] f.month[h.date]")
		fmt.Fprintln(w, "These can be used in any combination and repeated for the desired effect:")
		fmt.Fprintln(w, "\tfilter h.user-agent icontains .Kmail into summary h.user-agent f.year[h.date] calculate f.count filter c.year-date eq 2022 sort h.user-agent")
		fmt.Fprintln(w, "And so forth")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Notes:")
		fmt.Fprintln(w, "- Single word string literals begin with a `.`")
		fmt.Fprintln(w, "- Headers are referred to with `h.` once in table form it becomes a column referred to by `c.`")
		fmt.Fprintln(w, "- I haven't implemented any body parsing components yet")
		fmt.Fprintln(w, "- Extension PRs are welcome and intended")
		fmt.Fprintln(w, "- All functions are preceded by `f.`")
		fmt.Fprintln(w, "- Once converted to table form, it can not be converted back to mailfile or mbox")
	}
	fmt.Fprintln(w, "A complete list of functions supported:")
	// TODO funcs.PrintFunctionList(w) when updated
	funcs.PrintFunctionList()
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "List of supported input types:")
	PrintInputHelp(w)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "List of supported output types: (Must be supported based on query.)")
	// TODO dataformats.PrintOutputHelp(w, customOutputs) when updated
	dataformats.PrintOutputHelp(customOutputs)
	fmt.Fprintln(w, "")
}
