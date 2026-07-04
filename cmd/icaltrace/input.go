package main

import (
	"fmt"
	"io"
	_ "github.com/emersion/go-message/charset"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/icaldata"
)

func InputHandler(inputType string, inputFile string, w io.Writer) (pimtrace.Data, error) {
	if w == nil {
		w = os.Stdout
	}
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
		PrintInputHelp(w)
	default:
		return nil, fmt.Errorf("please specify an -input-type. got %s", inputType)
	}
	return icaldata.Data(ventry), nil
}

func PrintInputHelp(w io.Writer) {
	_, _ = fmt.Fprintln(w, "input-types available: ")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "ical", "Read an iCal file or '-' for stdin")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "list", "This help text")
	_, _ = fmt.Fprintln(w)
}
