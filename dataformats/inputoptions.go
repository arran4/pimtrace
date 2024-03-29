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
		var lastProgress = -1
		var lastDuration time.Duration = 0
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
			_, err = seeker.Seek(int64(position), io.SeekStart)
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
				now := time.Now()
				duration := now.Sub(start)
				log.Printf("Done %d of %d in %s", position, end, duration)
				return nil
			},
			NewRead: func(p []byte) (n int, err error) {
				n, err = reader.Read(p)
				if !ok {
					return
				}
				position += n
				pct := (position * 100) / end
				if pct != progress {
					progress = pct
					now := time.Now()
					var duration time.Duration
					var estimate time.Duration
					if progress > 0 {
						duration = now.Sub(start)
						estimate = (100 * duration) / time.Duration(progress)
					}
					if duration-lastDuration > 2*time.Second || progress-lastProgress > 25 {
						log.Printf("%d%%; bytes: %d/%d; duration: %s / %s", progress, position, end, duration, estimate)
						lastProgress = progress
						lastDuration = duration
					}
				}
				return
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
