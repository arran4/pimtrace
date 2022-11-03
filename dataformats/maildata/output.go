package maildata

import (
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message/textproto"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
)

func (p MailDataType) WriteCSVStream(stdin *os.File, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p MailDataType) WriteCSVFile(s string) error {
	//TODO implement me
	panic("implement me")
}

func (p MailDataType) WriteMBoxFile(fName string) error {
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return p.WriteMBoxStream(f, fName)
}

func (p MailDataType) WriteMBoxStream(f io.Writer, fName string) error {
	mbw := mbox.NewWriter(f)
	for mi, m := range p {
		mw, err := mbw.CreateMessage(m.From(), m.Time())
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("creating message %d to Mbox %s: %w", mi+1, fName, err)
		}
		if err := p.WriteMailStream(mw, fName); err != nil {
			return fmt.Errorf("writing message %d to Mbox %s: %w", mi+1, fName, err)
		}
	}
	return nil
}

func (p MailDataType) WriteMailFile(fName string) error {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("writing mail file %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return p.WriteMailStream(f, fName)
}

func (p MailDataType) WriteMailStream(f io.Writer, fName string) error {
	for mi, m := range p {
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
