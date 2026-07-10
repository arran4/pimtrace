package groupdata

import (
	"io"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
	"pimtrace/fsys"
)

func (d Data) WriteCSVFile(fName string) error {
	return pimtrace.WriteFileWrapper(fsys.OSFS{}, "CSV", fName, d.WriteCSVStream)
}

func (d Data) WriteCSVStream(f io.Writer, fName string) error {
	return tabledata.WriteCsv(d, f)
}

func (d Data) WriteTableFile(fName string) error {
	return pimtrace.WriteFileWrapper(fsys.OSFS{}, "Table", fName, d.WriteTableStream)
}

func (d Data) WriteTableStream(f io.Writer, fName string) error {
	tabledata.WriteTable(d, f)
	return nil
}
