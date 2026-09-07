package pimtrace_test

import (
	"bytes"
	"io"
	"pimtrace"
	"pimtrace/dataformats"
	"pimtrace/dataformats/tabledata"
	"testing"
)

type failWriter struct{}

func (f *failWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

type trackCloser struct {
	io.Writer
	closed bool
}

func (t *trackCloser) Close() error {
	t.closed = true
	return nil
}

func TestOutputHandler_ErrorsAndClosures(t *testing.T) {
	data := tabledata.Data{
		&tabledata.Row{
			Headers: map[string]int{"Header1": 0},
			Row:     []pimtrace.Value{pimtrace.SimpleStringValue("Value1")},
		},
	}

	t.Run("failing writer on count", func(t *testing.T) {
		err := dataformats.OutputHandler(data, "count", "-", nil, &failWriter{})
		if err == nil {
			t.Errorf("expected error from failing writer, got nil")
		}
	})

	t.Run("ensure streams are not closed", func(t *testing.T) {
		tc := &trackCloser{Writer: &bytes.Buffer{}}
		err := dataformats.OutputHandler(data, "csv", "-", nil, tc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tc.closed {
			t.Errorf("expected injected writer to remain open, but it was closed")
		}
	})
}
