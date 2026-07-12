package fsystest

import (
	"pimtrace/fsys"
	"io"
	"os"
	"testing/fstest"
)

type MapFSAdapter struct {
	MapFS fstest.MapFS
}

type NopFile struct {
	io.ReadCloser
}

func (n NopFile) Write(p []byte) (int, error) { return 0, io.ErrClosedPipe }

func (n NopFile) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := n.ReadCloser.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, io.ErrUnexpectedEOF
}

func (m MapFSAdapter) OpenFile(name string, flag int, perm os.FileMode) (fsys.File, error) {
	f, err := m.MapFS.Open(name)
	if err != nil {
		return nil, err
	}
	return NopFile{ReadCloser: f}, nil
}
