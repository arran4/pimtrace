package dataformats

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Next[T any] func(f io.Reader, fType string, fName string, ops ...any) ([]T, error)

func ReadFile[T any](fType string, fName string, next Next[T], ops ...any) ([]T, error) {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading %s %s: %w", fType, fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	ff, closers, err := ReaderStreamMapperOptionProcessor(f, ops)
	defer func() {
		for i := range closers {
			fc := closers[len(closers)-i-1]
			if err := fc.Close(); err != nil {
				log.Printf("error closing ReaderStreamMapper: %s", err)
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	return next(ff, fType, fName)
}
