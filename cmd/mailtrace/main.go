package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pimtrace/argparsers/basic"
	_ "pimtrace/funcs"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	f := flag.FlagSet{
		Usage: func() {
			log.Printf("Usage")
		},
	}
	var (
		inputType   = f.String("input-type", "list", "The input type")
		inputFile   = f.String("input", "-", "Input file or - for stdin")
		outputType  = f.String("output-type", "list", "The input type")
		outputFile  = f.String("output", "-", "Output file or - for stdin")
		versionFlag = f.Bool("version", false, "Prints the version")
	)

	if *versionFlag {
		fmt.Println(version)
		return
	}

	if err := f.Parse(os.Args); err != nil {
		log.Printf("Error parsing flags: %s", err)
		os.Exit(-1)
	}

	data, err := InputHandler(*inputType, *inputFile)
	if err != nil {
		log.Printf("Read Error: %s", err)
		os.Exit(-1)
	}

	ops, err := basic.ParseOperations(f.Args())
	if err != nil {
		log.Printf("Parse Error: %s", err)
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
