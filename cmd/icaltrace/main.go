package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pimtrace/argparsers/basic"
	"pimtrace/ast"
	"pimtrace/funcs"
	_ "pimtrace/funcs"
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
		PrintQueryHelp(*parser)
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

	data, err := InputHandler(*inputType, *inputFile)
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
		data, err = ops.Execute(data)
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

func PrintQueryHelp(parser string) {
	fmt.Println("This tool is for helping you filter, query and summarize ical files in a comprehensible way")
	fmt.Println("The usage is as follows:")
	fmt.Println("\ticaltrace -parser basic -input events.ical -input-type ical -output table $QUERY")
	fmt.Println("In this example it selects the basic parser, reads from events.ical, of the type ical. Outputs a table")
	fmt.Println("and runs query $QUERY. You are required to specify all of these arguments.")
	fmt.Println("")
	switch parser {
	case "basic":
		fmt.Println("Basic Parser")
		fmt.Println("")
		fmt.Println("Queries can do the following:")
		fmt.Println("- Filtering out data")
		fmt.Println("- Selecting components of the data to view")
		fmt.Println("- Grouping and summarizing data")
		fmt.Println("")
		fmt.Println("The Queries can be build up like this:")
		fmt.Println("Simple filter query")
		fmt.Println("\tfilter not p.SUMMARY icontains .Report")
		fmt.Println("Which filters all events without SUMMARY containing 'Report' ")
		fmt.Println("Then you can convert the emails to tabular form using:")
		fmt.Println("\tinto table p.LOCATION f.year[p.DUE] f.month[p.DUE]")
		fmt.Println("Which creates table of locations, and the due year and due month.")
		fmt.Println("If you wanted a summary / count of the lines you can use the summary converter:")
		fmt.Println("\tinto summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count")
		fmt.Println("Which groups the calendar invites based on location, due date year, due date month, then adds a count")
		fmt.Println("If you want to sort you can use:")
		fmt.Println("\tsort f.year[p.DUE] f.month[p.DUE]")
		fmt.Println("These can be used in any combination and repeated for the desired effect:")
		fmt.Println("\tfilter p.SUMMARY icontains .Report into summary p.LOCATION f.year[p.DUE] f.month[p.DUE] calculate f.count filter c.year-DUE eq 2022 sort p.LOCATION")
		fmt.Println("And so forth")
		fmt.Println("")
		fmt.Println("Notes:")
		fmt.Println("- Single word string literals begin with a `.`")
		fmt.Println("- Properties are referred to with `p.` once in table form it becomes a column referred to by `c.`")
		fmt.Println("- I haven't implemented sub components yet")
		fmt.Println("- Extension PRs are welcome and intended")
		fmt.Println("- All functions are preceded by `f.`")
		fmt.Println("- Once converted to table form, it can not be converted back to ical / ics")
	}
	fmt.Println("A complete list of functions supported:")
	funcs.PrintFunctionList()
	fmt.Println("")
	fmt.Println("List of supported input types:")
	PrintInputHelp()
	fmt.Println("")
	fmt.Println("List of supported output types: (Must be supported based on query.)")
	PrintOutputHelp()
	fmt.Println("")
}
