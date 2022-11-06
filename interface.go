package pimtrace

import "io"

type MailFileOutputCapable interface {
	WriteMailFile(fName string) error
	WriteMailStream(f io.Writer, fName string) error
}

type ICalFileOutputCapable interface {
	WriteICalFile(fName string) error
	WriteICalStream(f io.Writer, fName string) error
}

type MBoxOutputCapable interface {
	WriteMBoxFile(fName string) error
	WriteMBoxStream(f io.Writer, fName string) error
}

type CSVOutputCapable interface {
	WriteCSVFile(fName string) error
	WriteCSVStream(f io.Writer, fName string) error
}

type TableOutputCapable interface {
	WriteTableFile(fName string) error
	WriteTableStream(f io.Writer, fName string) error
}
