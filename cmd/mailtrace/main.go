package main

import (
	"flag"
	"log"
	"os"
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
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}

	ops := ParseOperations(flag.Args())
	if ops != nil {
		data, err = ops.Execute(data)
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(-1)
		}
	}

	if err := OutputHandler(data, outputType, outputFile); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}
}
