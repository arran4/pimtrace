package main

import (
	"fmt"
	"os"
	"pimtrace"
	"reflect"
)

func OutputHandler(p pimtrace.Data, mode, outputPath string) error {
	switch mode {
	case "ical":
		if np, ok := p.(pimtrace.ICalFileOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteICalStream(os.Stdin, outputPath)
			default:
				return np.WriteICalFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
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
	default:
		//fmt.Println("Please specify a -input-type")
		//fmt.Println()
		return nil
	}
}

func PrintOutputHelp() {
	fmt.Println("`--output-type`s: ")
	fmt.Printf(" %-30s %s\n", "ical", "iCal file as per: https://github.com/arran4/golang-ical")
	fmt.Printf(" %-30s %s\n", "list", "This help text")
	fmt.Printf(" %-30s %s\n", "csv", "Data in csv format")
	fmt.Printf(" %-30s %s\n", "count", "Just a count of rows")
	fmt.Println()
}
