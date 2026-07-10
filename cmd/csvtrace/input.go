package main

import (
	"fmt"
	"io"
	_ "github.com/emersion/go-message/charset"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/tabledata"
	"pimtrace/fsys"
)

func InputHandler(fs fsys.FS, inputType string, inputFile string, w io.Writer) (pimtrace.Data, error) {
	if w == nil {
		w = os.Stdout
	}
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
			nm, err := dataformats.ReadFile(fs, inputType, inputFile, tabledata.ReadCSV)
			if err != nil {
				return nil, err
			}
			rows = append(rows, nm...)
		}
	case "list":
		PrintInputHelp(w)
	default:
		return nil, fmt.Errorf("please specify an -input-type. got %s", inputType)
	}
	return tabledata.Data(rows), nil
}

func PrintInputHelp(w io.Writer) {
	_, _ = fmt.Fprintln(w, "input-types available: ")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "csv", "Read a CSV file")
	_, _ = fmt.Fprintf(w, " %-30s %s\n", "list", "This help text")
	_, _ = fmt.Fprintln(w)
}
