package dataformats

import (
	"fmt"
	"io"
	"log"
	"os"
	"pimtrace/fsys"
)

type Next[T any] func(f io.Reader, fType string, fName string, ops ...any) ([]T, error)

func ReadFile[T any](fType string, fName string, next Next[T], ops ...any) (res []T, err error) {
	fs := fsys.NewOSFS()
	for _, op := range ops {
		if o, ok := op.(fsys.FS); ok {
			fs = o
		}
	}
	f, err := fs.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading %s %s: %w", fType, fName, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			if err == nil {
				err = fmt.Errorf("closing file %s: %w", fName, cerr)
			} else {
				log.Printf("Error closing file: %s: %s", fName, cerr)
			}
		}
	}()
	ff, closers, err := ReaderStreamMapperOptionProcessor(f, ops)
	defer func() {
		for i := range closers {
			fc := closers[len(closers)-i-1]
			if cerr := fc.Close(); cerr != nil {
				if err == nil {
					err = fmt.Errorf("closing ReaderStreamMapper: %w", cerr)
				} else {
					log.Printf("error closing ReaderStreamMapper: %s", cerr)
				}
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	return next(ff, fType, fName)
}
