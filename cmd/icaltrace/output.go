package main

import (
	"fmt"
	"io"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"reflect"
)

var (
	customOutputs = [][2]string{
		{"ical", "iCal file as per: https://github.com/arran4/golang-ical"},
	}
)

func OutputHandler(p pimtrace.Data, mode, outputPath string, ops ...any) error {
	var out io.Writer = os.Stdout
	for _, op := range ops {
		if o, ok := op.(io.Writer); ok {
			out = o
		}
	}
	switch mode {
	case "ical":
		if np, ok := p.(pimtrace.ICalFileOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteICalStream(out, outputPath)
			default:
				return np.WriteICalFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	}
	return dataformats.OutputHandler(p, mode, outputPath, customOutputs, ops...)
}
