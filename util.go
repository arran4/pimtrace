package pimtrace

import (
	"fmt"
	"io"
	"log"
	"os"
)

type HasStringArray interface {
	StringArray(header []string) []string
	HeadersStringArray() []string
}

func WriteFileWrapper(fType string, fName string, fun func(f io.Writer, fName string) error) error {
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating %s %s: %w", fType, fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing %s file: %s: %s", fType, fName, err)
		}
	}()
	return fun(f, fName)
}
