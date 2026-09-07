package dataformats

import (
	"fmt"
	"io"
	"os"
	"pimtrace"
	"pimtrace/dataformats/plotoutput"
	"reflect"
)

func OutputHandler(p pimtrace.Data, mode, outputPath string, customOutputs [][2]string, ops ...any) error {
	var out io.Writer = os.Stdout
	for _, op := range ops {
		if o, ok := op.(io.Writer); ok {
			out = o
		}
	}
	switch mode {
	case "csv":
		if np, ok := p.(pimtrace.CSVOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteCSVStream(out, outputPath)
			default:
				return np.WriteCSVFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "table":
		if np, ok := p.(pimtrace.TableOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteTableStream(out, outputPath)
			default:
				return np.WriteTableFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "count":
		_, _ = fmt.Fprintln(out, p.Len())
		return nil
	case "list":
		PrintOutputHelp(out, customOutputs)
		return nil
	case "plot.bar":
		if outputPath == "-" {
			return fmt.Errorf("plot requires an -output file name rather than: `-output=%s`", outputPath)
		}
		return plotoutput.BarPlot(p, outputPath)
	default:
		return fmt.Errorf("please specify an -output-type")
	}
}

func PrintOutputHelp(out io.Writer, custom [][2]string) {
	_, _ = fmt.Fprintln(out, "--output-types: ")
	each := [][2]string{
		{"list", "This help text"},
		{"csv", "Data in csv format"},
		{"table", "Data in a ascii table"},
		{"count", "Just a count of rows"},
		{"plot.bar", "Writes a plot of the data out, the data must be tabular and columns must be in the form of: string, number*"},
	}
	for _, e := range append(each, custom...) {
		_, _ = fmt.Fprintf(out, " %-30s %s\n", e[0], e[1])
	}
	_, _ = fmt.Fprintln(out)
}
