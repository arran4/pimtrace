package dataformats

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
)

type ReaderStreamMapper func(io.Reader) (io.Reader, error)

var Gzip ReaderStreamMapper = fgzip

func fgzip(reader io.Reader) (io.Reader, error) {
	return gzip.NewReader(reader)
}

var Bzip2 ReaderStreamMapper = fbzip2

func fbzip2(reader io.Reader) (io.Reader, error) {
	return bzip2.NewReader(reader), nil
}

func ReaderStreamMapperOptionProcessor(f io.Reader, ops []any) (io.Reader, []io.Closer, error) {
	var ff = f
	var closers []io.Closer
	for i, op := range ops {
		switch op := op.(type) {
		case ReaderStreamMapper:
			var err error
			ff, err = op(f)
			if err != nil {
				return nil, nil, fmt.Errorf("with ReaderStreamMapper option %d: %w", i, err)
			}
			if fc, ok := ff.(io.Closer); ok {
				closers = append(closers, fc)
			}
		default:
			return nil, closers, fmt.Errorf("unknown option: %d", i)
		}
	}
	return ff, closers, nil
}
