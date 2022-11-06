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
	"strings"
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

	if err := f.Parse(os.Args); err != nil {
		log.Printf("Error parsing flags: %s", err)
		os.Exit(-1)
	}

	if f.NArg() <= 1 || *helpFlag {
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
	fmt.Println("This tool is for helping you filter, query and summarize CSV/data files in a comprehensible way")
	fmt.Println("The usage is as follows:")
	fmt.Println("\tcsvtrace -parser basic -input mail.mbox -input-type csv -output table $QUERY")
	fmt.Println("In this example it selects the basic parser (there is only one I intend to extend it if I get time")
	fmt.Println("but to avoid issues when I change it I am requiring specification.")
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
		fmt.Println("\tfilter not c.job icontains .Software")
		fmt.Println("Which filters all row's `job` column that does not contain `Software` ")
		fmt.Println("Then you can filter the columns using:")
		fmt.Println("\tinto table c.job f.year[c.startdate] f.month[c.startdate]")
		fmt.Println("If you wanted a summary / count of the lines you can use the summary converter:")
		fmt.Println("\tinto summary c.job calculate f.sum[c.pay] f.count")
		fmt.Println("Which groups the row on job, then creates a sum of pay and a count of elements grouped")
		fmt.Println("If you want to sort you can use:")
		fmt.Println("\tsort f.year[c.startdate] f.month[c.startdate]")
		fmt.Println("These can be used in any combination and repeated for the desired effect:")
		fmt.Println("\tfilter c.startdate icontains .Software into summary c.job f.year[c.startdate] calculate f.count filter c.year-startdate eq 2022 sort c.job")
		fmt.Println("And so forth")
		fmt.Println("")
		fmt.Println("Notes:")
		fmt.Println("- Single word string literals begin with a `.`")
		fmt.Println("- Columns are referred to by `c.`")
		fmt.Println("- Extension PRs are welcome and intended")
		fmt.Println("- All functions are preceded by `f.`")
	}
	fmt.Println("A complete list of functions supported:")
	PrintFunctionList()
	fmt.Println("")
	fmt.Println("List of supported input types:")
	PrintInputHelp()
	fmt.Println("")
	fmt.Println("List of supported output types: (Must be supported based on query.)")
	PrintOutputHelp()
	fmt.Println("")
}

func PrintFunctionList() {
	fmt.Println("Functions: ")
	for _, f := range funcs.Functions[ast.ValueExpression]() {
		for _, af := range f.Arguments() {
			args := make([]string, 0, len(af.Args))
			for _, aff := range af.Args {
				args = append(args, aff.String())
			}
			fn := fmt.Sprintf("f.%s[%s]", f.Name(), strings.Join(args, ","))
			fmt.Printf("%-40s%40s\n", fn, af.Description)
		}
	}
}
