package groupdata

import (
	"io"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
)

func (d Data) WriteCSVFile(fName string) error {
	return pimtrace.WriteFileWrapper("CSV", fName, d.WriteCSVStream)
}

func (d Data) WriteCSVStream(f io.Writer, fName string) error {
	return tabledata.WriteCsv(d, f)
}

func (d Data) WriteTableFile(fName string) error {
	return pimtrace.WriteFileWrapper("Table", fName, d.WriteTableStream)
}

func (d Data) WriteTableStream(f io.Writer, fName string) error {
	tabledata.WriteTable(d, f)
	return nil
}
