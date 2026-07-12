package dataformats

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"pimtrace/fsys"
)

func ReadTarFile[T any](fType string, fName string, next Next[T], globs []string, ops ...any) (res []T, err error) {
	var fs fsys.FS = fsys.NewOSFS()
	for _, op := range ops {
		if o, ok := op.(fsys.FS); ok {
			fs = o
		}
	}
	f, err := fs.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading Mbox %s: %w", fName, err)
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
	return ReadTarStream(ff, fType, fName, next, globs)
}

func ReadTarStream[T any](f io.Reader, fType string, fName string, next Next[T], globs []string, ops ...any) (res []T, err error) {
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
	t := tar.NewReader(ff)
	var ta []T
	for {
		ht, err := t.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tar: %w", err)
		}
		_, fFilename := filepath.Split(ht.Name)
		m := false
		for _, g := range globs {
			if ok, err := filepath.Match(g, fFilename); ok {
				m = true
			} else if err != nil {
				return nil, fmt.Errorf("filepath match %#v: %w", g, err)
			}
		}
		if !m {
			continue
		}
		taa, err := next(t, fType, ht.Name)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("tar: %w", err)
		}
		ta = append(ta, taa...)
	}
	return ta, nil
}
