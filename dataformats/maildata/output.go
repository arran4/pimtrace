package maildata

import (
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message/textproto"
	"io"
	"mime"
	"mime/multipart"
	"pimtrace"
	"pimtrace/dataformats/tabledata"
)

var _ pimtrace.MailFileOutputCapable = (*Data)(nil)
var _ pimtrace.MBoxOutputCapable = (*Data)(nil)
var _ pimtrace.CSVOutputCapable = (*Data)(nil)
var _ pimtrace.TableOutputCapable = (*Data)(nil)

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

func (mdt Data) WriteMBoxFile(fName string) error {
	return pimtrace.WriteFileWrapper("MBox", fName, mdt.WriteMBoxStream)
}

func (mdt Data) WriteMBoxStream(f io.Writer, fName string) error {
	mbw := mbox.NewWriter(f)
	for mi, m := range mdt {
		mw, err := mbw.CreateMessage(m.From(), m.Time())
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("creating message %d to Mbox %s: %w", mi+1, fName, err)
		}
		if err := mdt.WriteMailStream(mw, fName); err != nil {
			return fmt.Errorf("writing message %d to Mbox %s: %w", mi+1, fName, err)
		}
	}
	return nil
}

func (mdt Data) WriteMailFile(fName string) error {
	return pimtrace.WriteFileWrapper("MailFile", fName, mdt.WriteMailStream)
}

func (mdt Data) WriteMailStream(f io.Writer, fName string) error {
	for mi, m := range mdt {
		if err := textproto.WriteHeader(f, textproto.HeaderFromMap(m.MailHeader.Map())); err != nil {
			return fmt.Errorf("writing message %d header %s: %w", mi+1, fName, err)
		}
		ct := m.MailHeader.Get("Content-Type")
		mt, mtp, err := mime.ParseMediaType(ct)
		if err != nil && ct != "" {
			return fmt.Errorf("writing message %d header %s content type : %w", mi+1, fName, err)
		}
		switch mt {
		case "multipart/alternative":
			mpw := multipart.NewWriter(f)
			if err := mpw.SetBoundary(mtp["Boundary"]); err != nil {
				return fmt.Errorf("multipart boundary error message %d header %s content type: %w", mi+1, fName, err)
			}
			for _, mb := range m.MailBodies {
				mpwf, err := mpw.CreatePart(mb.Header())
				if err != nil {
					return fmt.Errorf("creating part for multipart message %d body %s: %w", mi+1, fName, err)
				}
				if _, err := io.Copy(mpwf, mb.Reader()); err != nil && !errors.Is(err, io.EOF) {
					return fmt.Errorf("writing message %d body %s: %w", mi+1, fName, err)
				}
			}
		default:
			for _, mb := range m.MailBodies {
				if _, err := io.Copy(f, mb.Reader()); err != nil && !errors.Is(err, io.EOF) {
					return fmt.Errorf("writing message %d body %s: %w", mi+1, fName, err)
				}
			}
		}
	}
	return nil
}
