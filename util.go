package pimtrace

import (
	"fmt"
	"io"
	"log"
	"os"
	"pimtrace/fsys"
)

type HasStringArray interface {
	StringArray(header []string) []string
	HeadersStringArray() []string
}

func WriteFileWrapper(fs fsys.FS, fType string, fName string, fun func(f io.Writer, fName string) error) (err error) {
	f, err := fs.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating %s %s: %w", fType, fName, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			if err == nil {
				err = fmt.Errorf("closing %s file %s: %w", fType, fName, cerr)
			} else {
				log.Printf("Error closing %s file: %s: %s", fType, fName, cerr)
			}
		}
	}()
	return fun(f, fName)
}

func ReadFileWrapper[T Data](fs fsys.FS, fType string, fName string, fun func(f io.Reader, fName string) (T, error)) (result T, err error) {
	f, err := fs.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return result, fmt.Errorf("reading %s %s: %w", fType, fName, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			if err == nil {
				err = fmt.Errorf("closing %s file %s: %w", fType, fName, cerr)
			} else {
				log.Printf("Error closing %s file: %s: %s", fType, fName, cerr)
			}
		}
	}()
	return fun(f, fName)
}
