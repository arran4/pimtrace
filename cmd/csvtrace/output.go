package main

import (
	"fmt"
	"os"
	"pimtrace"
	"pimtrace/dataformats/plotoutput"
	"reflect"
)

func OutputHandler(p pimtrace.Data, mode, outputPath string) error {
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
		PrintOutputHelp()
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

func PrintOutputHelp() {
	fmt.Println("`--output-type`s: ")
	fmt.Printf(" %-30s %s\n", "list", "This help text")
	fmt.Printf(" %-30s %s\n", "csv", "Data in csv format")
	fmt.Printf(" %-30s %s\n", "table", "Data in a ascii table")
	fmt.Printf(" %-30s %s\n", "count", "Just a count of rows")
	fmt.Printf(" %-30s %s\n", "plot.bar", "Writes a plot of the data out, the data must be tabular and columns must be in the form of: string, number*")
	fmt.Println()
}
