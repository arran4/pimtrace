package main

import (
	"fmt"
	"os"
	"pimtrace"
	"pimtrace/dataformats"
	"reflect"
)

var (
	customOutputs = [][2]string{
		{"mailfile", "A single mail file (do not use see https://groups.google.com/g/golang-nuts/c/T1xoNVr6ask/m/gdtdRUShCwAJ)"},
		{"mbox", "Mbox file (do not use see https://groups.google.com/g/golang-nuts/c/T1xoNVr6ask/m/gdtdRUShCwAJ)"},
	}
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
	}
	return dataformats.OutputHandler(p, mode, outputPath, customOutputs)
}
