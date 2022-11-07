package main

import (
	"fmt"
	"os"
	"pimtrace"
	"reflect"
)

func OutputHandler(p pimtrace.Data, mode, outputPath string) error {
	switch mode {
	case "mailfile":
		if np, ok := p.(pimtrace.MailFileOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteMailStream(os.Stdin, outputPath)
			default:
				return np.WriteMailFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "mbox":
		if np, ok := p.(pimtrace.MBoxOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteMBoxStream(os.Stdin, outputPath)
			default:
				return np.WriteMBoxFile(outputPath)
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
		return fmt.Errorf("please specify an -output-type")
	}
}

func PrintOutputHelp() {
	fmt.Println("`--output-type`s: ")
	fmt.Printf(" %-30s %s\n", "mailfile", "A single mail file (do not use see https://groups.google.com/g/golang-nuts/c/T1xoNVr6ask/m/gdtdRUShCwAJ)")
	fmt.Printf(" %-30s %s\n", "mbox", "Mbox file (do not use see https://groups.google.com/g/golang-nuts/c/T1xoNVr6ask/m/gdtdRUShCwAJ)")
	fmt.Printf(" %-30s %s\n", "list", "This help text")
	fmt.Printf(" %-30s %s\n", "csv", "Data in csv format")
	fmt.Printf(" %-30s %s\n", "table", "Data in a ascii table")
	fmt.Printf(" %-30s %s\n", "count", "Just a count of rows")
	fmt.Println()
}
