package dataformats

import (
	"fmt"
	"os"
	"pimtrace"
	"pimtrace/dataformats/plotoutput"
	"reflect"
)

func OutputHandler(p pimtrace.Data, mode, outputPath string, customOutputs [][2]string) error {
	switch mode {
	case "csv":
		if np, ok := p.(pimtrace.CSVOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteCSVStream(os.Stdin, outputPath)
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
				return np.WriteTableStream(os.Stdin, outputPath)
			default:
				return np.WriteTableFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "count":
		fmt.Println(p.Len())
		return nil
	case "list":
		PrintOutputHelp(customOutputs)
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

func PrintOutputHelp(custom [][2]string) {
	fmt.Println("`--output-type`s: ")
	each := [][2]string{
		{"list", "This help text"},
		{"csv", "Data in csv format"},
		{"table", "Data in a ascii table"},
		{"count", "Just a count of rows"},
		{"plot.bar", "Writes a plot of the data out, the data must be tabular and columns must be in the form of: string, number*"},
	}
	for _, e := range append(each, custom...) {
		fmt.Printf(" %-30s %s\n", e[0], e[1])
	}
	fmt.Println()
}
