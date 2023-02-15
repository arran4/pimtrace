package dataformats

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"time"
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

type CloseOverwriter struct {
	NewClose func() error
	NewRead  func(p []byte) (n int, err error)
}

func (c *CloseOverwriter) Read(p []byte) (n int, err error) {
	return c.NewRead(p)
}

func (c *CloseOverwriter) Close() error {
	return c.NewClose()
}

var _ io.ReadCloser = (*CloseOverwriter)(nil)

func NewProgressor() ReaderStreamMapper {
	return func(reader io.Reader) (io.Reader, error) {
		seeker, ok := reader.(io.Seeker)
		var position = -1
		var progress = -1
		var end = -1
		start := time.Now()
		if ok {
			n, err := seeker.Seek(0, io.SeekCurrent)
			if err != nil {
				log.Printf("Error with Progressor: %s", err)
				ok = false
			}
			position = int(n)
			n, err = seeker.Seek(0, io.SeekEnd)
			if err != nil {
				log.Printf("Error with Progressor: %s", err)
				ok = false
			}
			end = int(n)
			n, err = seeker.Seek(int64(position), io.SeekStart)
			if err != nil {
				log.Printf("Error with Progressor: %s", err)
				ok = false
			}
		}
		return &CloseOverwriter{
			NewClose: func() error {
				if !ok {
					return nil
				}
				log.Printf("Done %d of %d", position, end)
				return nil
			},
			NewRead: func(p []byte) (n int, err error) {
				if !ok {
					return reader.Read(p)
				}
				position += n
				pct := position * 100 / end
				if pct != progress {
					progress = pct
					now := time.Now()
					var duration time.Duration
					var estimate time.Duration
					if progress > 0 {
						duration = now.Sub(start)
						estimate = (100 * duration) / time.Duration(progress)
					}
					log.Printf("%d%% bytes: %d/%d duration: %s / %s", progress, position, end, duration, estimate)
				}
				return reader.Read(p)
			},
		}, nil
	}
}

func ReaderStreamMapperOptionProcessor(f io.Reader, ops []any) (io.Reader, []io.Closer, error) {
	var ff = f
	var closers []io.Closer
	for i, op := range ops {
		switch op := op.(type) {
		case ReaderStreamMapper:
			var err error
			ff, err = op(ff)
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
