package groupdata

import (
	"io"
	"pimtrace"
)

var _ pimtrace.CSVOutputCapable = (*Data)(nil)

func (d Data) WriteCSVFile(fName string) error {
	//TODO implement me
	panic("implement me")
}

func (d Data) WriteCSVStream(f io.Writer, fName string) error {
	//TODO implement me
	panic("implement me")
}

var _ pimtrace.TableOutputCapable = (*Data)(nil)

func (d Data) WriteTableFile(fName string) error {
	//TODO implement me
	panic("implement me")
}

func (d Data) WriteTableStream(f io.Writer, fName string) error {
	//TODO implement me
	panic("implement me")
}
