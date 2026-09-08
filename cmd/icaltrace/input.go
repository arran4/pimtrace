package main

import (
	"fmt"
	_ "github.com/emersion/go-message/charset"
	"io"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/icaldata"
)

func InputHandler(inputType string, inputFile string, ops ...any) (pimtrace.Data, error) {
	var w io.Writer = os.Stdout
	var inputReader io.Reader = os.Stdin
	for _, op := range ops {
		if o, ok := op.(io.Writer); ok {
			w = o
		}
		if in, ok := op.(io.Reader); ok {
			inputReader = in
		}
	}
	ventry := []*icaldata.ICalWithSource{}
	switch inputType {
	case "ical":
		switch inputFile {
		case "-":
			nm, err := icaldata.ReadICalStream(inputReader, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			ventry = append(ventry, nm...)
		default:
			nm, err := dataformats.ReadFile(inputType, inputFile, icaldata.ReadICalStream, ops...)
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
