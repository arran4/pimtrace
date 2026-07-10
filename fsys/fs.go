package fsys

import (
	"io"
	"os"
)

type File interface {
	io.ReadWriteCloser
	io.Seeker
}

type FS interface {
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
}

type OSFS struct{}

func (OSFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}
