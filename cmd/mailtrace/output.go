package main

import (
	"fmt"
	"io"
	"os"
	"pimtrace"
	"reflect"
)

type MailFileOutputCapable interface {
	WriteMailFile(fName string) error
	WriteMailStream(f io.Writer, fName string) error
}

type MBoxOutputCapable interface {
	WriteMBoxFile(fName string) error
	WriteMBoxStream(f io.Writer, fName string) error
}

type CSVOutputCapable interface {
	WriteCSVFile(fName string) error
	WriteCSVStream(f io.Writer, fName string) error
}

func OutputHandler(p pimtrace.Data, mode, outputPath string) error {
	switch mode {
	case "mailfile":
		if np, ok := p.(MailFileOutputCapable); ok {
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
		if np, ok := p.(MBoxOutputCapable); ok {
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
		if np, ok := p.(CSVOutputCapable); ok {
			switch outputPath {
			case "-":
				return np.WriteCSVStream(os.Stdin, outputPath)
			default:
				return np.WriteCSVFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "count":
		fmt.Println(p.Len())
		return nil
	case "list":
		fmt.Println("`--output-type`s: ")
		fmt.Printf(" %-30s %s\n", "mailfile", "A single mail file")
		fmt.Printf(" %-30s %s\n", "mbox", "Mbox file")
		fmt.Printf(" %-30s %s\n", "list", "This help text")
		fmt.Printf(" %-30s %s\n", "count", "Just a count")
		fmt.Printf(" %-30s %s\n", "csv", "Data in csv format")
		fmt.Println()
		return nil
	default:
		//fmt.Println("Please specify a -input-type")
		//fmt.Println()
		return nil
	}
}
