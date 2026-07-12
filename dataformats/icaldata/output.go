package icaldata

import (
	"fmt"
	ics "github.com/arran4/golang-ical"
	"io"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
)

var _ pimtrace.ICalFileOutputCapable = (*Data)(nil)
var _ pimtrace.CSVOutputCapable = (*Data)(nil)
var _ pimtrace.TableOutputCapable = (*Data)(nil)

func (icd Data) WriteCSVFile(fName string) error {
	return pimtrace.WriteFileWrapper("CSV", fName, icd.WriteCSVStream)
}

func (icd Data) WriteCSVStream(f io.Writer, fName string) error {
	return tabledata.WriteCsv(icd, f)
}

func (icd Data) WriteTableFile(fName string) error {
	return pimtrace.WriteFileWrapper("Table", fName, icd.WriteTableStream)
}

func (icd Data) WriteTableStream(f io.Writer, fName string) error {
	tabledata.WriteTable(icd, f)
	return nil
}

func (icd Data) WriteICalFile(fName string) error {
	return pimtrace.WriteFileWrapper("ICalFile", fName, icd.WriteICalStream)
}

func (icd Data) WriteICalStream(f io.Writer, fName string) error {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for _, m := range icd {
		cal.Components = append(cal.Components, m.Component)
	}
	err := cal.SerializeTo(f)
	if err != nil {
		return fmt.Errorf("write ical: %w", err)
	}
	return nil
}
