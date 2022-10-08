package main

import (
	"flag"
	"log"
	"os"
)

var (
	inputType = flag.String("input-type", "list", "The input type")
	inputFile = flag.String("input", "-", "Input file or - for stdin")
	//outputType = flag.String("output-type", "list", "The input type")
	outputFile = flag.String("output", "-", "Output file or - for stdin")
)

func main() {
	flag.Parse()
	mails, err := InputHandler()
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}

	if err := OutputHandler(mails); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}
}
