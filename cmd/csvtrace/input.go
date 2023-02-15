package main

import (
	"fmt"
	_ "github.com/emersion/go-message/charset"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/tabledata"
)

func InputHandler(inputType string, inputFile string) (pimtrace.Data, error) {
	var rows []*tabledata.Row
	switch inputType {
	case "csv":
		switch inputFile {
		case "-":
			nm, err := tabledata.ReadCSV(os.Stdin, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			rows = append(rows, nm...)
		default:
			nm, err := dataformats.ReadFile(inputType, inputFile, tabledata.ReadCSV)
			if err != nil {
				return nil, err
			}
			rows = append(rows, nm...)
		}
	case "list":
		PrintInputHelp()
	default:
		return nil, fmt.Errorf("please specify an -input-type. got %s", inputType)
	}
	return tabledata.Data(rows), nil
}

func PrintInputHelp() {
	fmt.Println("input-types available: ")
	fmt.Printf(" %-30s %s\n", "csv", "Read a CSV file")
	fmt.Printf(" %-30s %s\n", "list", "This help text")
	fmt.Println()
}
