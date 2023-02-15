package main

import (
	"fmt"
	_ "github.com/emersion/go-message/charset"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/icaldata"
)

func InputHandler(inputType string, inputFile string) (pimtrace.Data, error) {
	ventry := []*icaldata.ICalWithSource{}
	switch inputType {
	case "ical":
		switch inputFile {
		case "-":
			nm, err := icaldata.ReadICalStream(os.Stdin, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			ventry = append(ventry, nm...)
		default:
			nm, err := dataformats.ReadFile(inputType, inputFile, icaldata.ReadICalStream)
			if err != nil {
				return nil, err
			}
			ventry = append(ventry, nm...)
		}
	case "list":
		PrintInputHelp()
	default:
		return nil, fmt.Errorf("please specify an -input-type. got %s", inputType)
	}
	return icaldata.Data(ventry), nil
}

func PrintInputHelp() {
	fmt.Println("input-types available: ")
	fmt.Printf(" %-30s %s\n", "mailfile", "A single mail file")
	fmt.Printf(" %-30s %s\n", "mbox", "Mbox file")
	fmt.Printf(" %-30s %s\n", "list", "This help text")
	fmt.Println()
}
