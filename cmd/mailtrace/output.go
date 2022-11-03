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

func OutputHandler[T any](p pimtrace.Data[T], mode, outputPath string) error {
	switch mode {
	case "mailfile":
		if p, ok := p.(MailFileOutputCapable); ok {
			switch outputPath {
			case "-":
				return p.WriteMailStream(os.Stdin, outputPath)
			default:
				return p.WriteMailFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "mbox":
		if p, ok := p.(MBoxOutputCapable); ok {
			switch outputPath {
			case "-":
				return p.WriteMBoxStream(os.Stdin, outputPath)
			default:
				return p.WriteMBoxFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "csv":
		if p, ok := p.(CSVOutputCapable); ok {
			switch outputPath {
			case "-":
				return p.WriteCSVStream(os.Stdin, outputPath)
			default:
				return p.WriteCSVFile(outputPath)
			}
		} else {
			return fmt.Errorf("unsupported format: %s of %s", mode, reflect.TypeOf(p))
		}
	case "count":
		fmt.Println(p.Len())
		return nil
	case "list":
		fmt.Println("`--output-type`s: ")
		fmt.Printf(" =%-20s - %s\n", "mailfile", "A single mail file")
		fmt.Printf(" =%-20s - %s\n", "mbox", "Mbox file")
		fmt.Printf(" =%-20s - %s\n", "list", "This help text")
		fmt.Printf(" =%-20s - %s\n", "count", "Just a count")
		fmt.Printf(" =%-20s - %s\n", "csv", "Data in csv format")
		fmt.Println()
		return nil
	default:
		//fmt.Println("Please specify a -input-type")
		//fmt.Println()
		return nil
	}
}
