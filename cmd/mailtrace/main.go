package main

import (
	"flag"
	"log"
	"os"
	"pimtrace/argparsers/basic"
	_ "pimtrace/funcs"
	"pimtrace/internal/maildata"
)

func main() {
	var (
		inputType  = flag.String("input-type", "list", "The input type")
		inputFile  = flag.String("input", "-", "Input file or - for stdin")
		outputType = flag.String("output-type", "list", "The input type")
		outputFile = flag.String("output", "-", "Output file or - for stdin")
	)

	flag.Parse()

	data, err := InputHandler(*inputType, *inputFile)
	if err != nil {
		log.Printf("Read Error: %s", err)
		os.Exit(-1)
	}

	ops, err := basic.ParseOperations[*maildata.MailWithSource](flag.Args())
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

	if err := OutputHandler(data, outputType, outputFile); err != nil {
		log.Printf("Write Error: %s", err)
		os.Exit(-1)
	}
}
