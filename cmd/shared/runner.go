package shared

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"pimtrace"
	"pimtrace/argparsers/basic"
	"pimtrace/ast"
	"pimtrace/funcs"

	"github.com/arran4/go-evaluator"
)

type Config struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	Args []string

	Name string

	PrintQueryHelp func(w io.Writer, parser string)
	PrintVersion   func(w io.Writer)
	InputHandler   func(inputType string, inputFile string, ops ...any) (pimtrace.Data, error)
	OutputHandler  func(data pimtrace.Data, outputType string, outputFile string, stdout io.Writer) error

	Progressor bool
}

func Run(c *Config) int {
	if c.Stdout == nil {
		c.Stdout = os.Stdout
	}
	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}
	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}

	f := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	f.SetOutput(c.Stderr)

	var (
		inputType   = f.String("input-type", "list", "The input type")
		inputFile   = f.String("input", "-", "Input file or - for stdin")
		outputType  = f.String("output-type", "list", "The input type")
		outputFile  = f.String("output", "-", "Output file or - for stdout")
		parser      = f.String("parser", "", "Just use `basic`")
		versionFlag = f.Bool("version", false, "Prints the version")
		helpFlag    = f.Bool("help", false, "Prints help")
		progress    *bool
	)

	if c.Progressor {
		progress = f.Bool("progress", false, "Report progress")
	}

	f.Usage = func() {
		_, _ = fmt.Fprintln(c.Stderr, "Usage: ", c.Name, "[Flags]", "[Query]")
		f.PrintDefaults()
		if c.PrintQueryHelp != nil {
			c.PrintQueryHelp(c.Stdout, *parser)
		}
	}

	if err := f.Parse(c.Args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		log.Printf("Error parsing flags: %s", err)
		return 2
	}

	if *versionFlag {
		if c.PrintVersion != nil {
			c.PrintVersion(c.Stdout)
		}
		return 0
	}

	if *helpFlag {
		f.Usage()
		return 0
	}

	if len(c.Args) == 0 {
		_, _ = fmt.Fprintln(c.Stderr, "No query found")
		f.Usage()
		return 2
	}

	// if there are no un-parsed trailing arguments AND we are not requesting list help
	// Wait, in previous revision `c.Args` had the query missing error only when len(c.Args) <= 1, which let flags pass even without query.
	// But the user test explicitly runs: csvtrace -parser basic -input - -input-type csv -output - -output-type csv
	// Which has NO trailing arguments (f.NArg() == 0).
	// So we shouldn't fail if there's no query provided if it's explicitly valid or if we aren't strict about it here.
	// Let's drop the second query check and rely on parser/InputHandler/OutputHandler logic to succeed on empty queries.
	if f.NArg() == 0 && (*inputType == "list" || *outputType == "list") {
		// Just listing input/output formats is completely fine without a query
	} else if f.NArg() == 0 && *parser != "basic" && *parser != "" {
		// If they chose a parser but provided no query args? Actually basic parser can parse an empty query into a no-op!
	}

	var iops []any
	iops = append(iops, c.Stdin) // Ensure stdin injected

	if c.Progressor && progress != nil && *progress {
		iops = append(iops, "progressor")
	}

	data, err := c.InputHandler(*inputType, *inputFile, iops...)
	if err != nil {
		_, _ = fmt.Fprintf(c.Stderr, "Read Error: %s\n", err)
		return 1
	}

	var ops ast.Operation
	switch *parser {
	case "basic":
		ops, err = basic.ParseOperations(f.Args())
		if err != nil {
			_, _ = fmt.Fprintf(c.Stderr, "Parse Error: %s\n", err)
			return 2
		}
	default:
		_, _ = fmt.Fprintln(c.Stderr, "Please use -parser=basic parameter, as maybe one day a more advanced parser will be created")
		return 2
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
			_, _ = fmt.Fprintf(c.Stderr, "Execute Error: %s\n", err)
			return 1
		}
	}

	if c.OutputHandler != nil {
		if err := c.OutputHandler(data, *outputType, *outputFile, c.Stdout); err != nil {
			_, _ = fmt.Fprintf(c.Stderr, "Write Error: %s\n", err)
			return 1
		}
	}

	return 0
}
