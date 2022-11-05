package tabledata

import (
	"encoding/csv"
	"github.com/olekukonko/tablewriter"
	"io"
	"pimtrace"
)

var _ pimtrace.CSVOutputCapable = (*Data)(nil)

func (d Data) WriteCSVFile(fName string) error {
	return pimtrace.WriteFileWrapper("CSV", fName, d.WriteCSVStream)
}

func (d Data) WriteCSVStream(f io.Writer, fName string) error {
	return WriteCsv(d, f)
}

var _ pimtrace.TableOutputCapable = (*Data)(nil)

func (d Data) WriteTableFile(fName string) error {
	return pimtrace.WriteFileWrapper("Table", fName, d.WriteTableStream)
}

func (d Data) WriteTableStream(f io.Writer, fName string) error {
	WriteTable(d, f)
	return nil
}

func WriteTable[T pimtrace.HasStringArray](d []T, f io.Writer) {
	table := tablewriter.NewWriter(f)
	var headers []string
	for i, v := range d {
		if i == 0 {
			headers = v.HeadersStringArray()
			table.SetHeader(headers)
		}
		table.Append(v.StringArray(headers))
	}
	table.Render() // Send output
}

func WriteCsv[T pimtrace.HasStringArray](d []T, f io.Writer) error {
	table := csv.NewWriter(f)
	var headers []string
	for i, v := range d {
		if i == 0 {
			headers = v.HeadersStringArray()
			if err := table.Write(headers); err != nil {
				return err
			}
		}
		if err := table.Write(v.StringArray(headers)); err != nil {
			return err
		}
	}
	table.Flush()
	return table.Error()
}
