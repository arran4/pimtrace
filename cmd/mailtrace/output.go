package main

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

func OutputHandler(mails []*MailWithSource, outputType *string, outputFile *string) error {
	switch *outputType {
	case "mailfile":
		switch *outputFile {
		case "-":
			err := WriteMailStream(mails, os.Stdin, *outputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteMailFile(mails, *outputFile)
			if err != nil {
				return err
			}
		}
	case "mbox":
		switch *outputFile {
		case "-":
			err := WriteMBoxStream(mails, os.Stdin, *outputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteMBoxFile(mails, *outputFile)
			if err != nil {
				return err
			}
		}
	case "count":
		fmt.Println(len(mails))
	case "list":
		fmt.Println("`--output-type`s: ")
		fmt.Printf(" =%-20s - %s\n", "mailfile", "A single mail file")
		fmt.Printf(" =%-20s - %s\n", "mbox", "Mbox file")
		fmt.Printf(" =%-20s - %s\n", "list", "This help text")
		fmt.Printf(" =%-20s - %s\n", "count", "Just a count")
		fmt.Println()
	default:
		fmt.Println("Please specify a -input-type")
		fmt.Println()
	}
	return nil
}

func WriteMBoxFile(ms []*MailWithSource, fName string) error {
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return WriteMBoxStream(ms, f, fName)
}

func WriteMBoxStream(ms []*MailWithSource, f io.Writer, fName string) error {
	mbw := mbox.NewWriter(f)
	for mi, m := range ms {
		mw, err := mbw.CreateMessage(m.From(), m.Time())
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("creating message %d to Mbox %s: %w", mi+1, fName, err)
		}
		if err := WriteMailStream(ms[mi:mi+1], mw, fName); err != nil {
			return fmt.Errorf("writing message %d to Mbox %s: %w", mi+1, fName, err)
		}
	}
	return nil
}

func WriteMailFile(ms []*MailWithSource, fName string) error {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("reading Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return WriteMailStream(ms, f, fName)
}

func WriteMailStream(ms []*MailWithSource, f io.Writer, fName string) error {
	for _, m := range ms {
		if err := textproto.WriteHeader(f, textproto.HeaderFromMap(m.MailHeader.Map())); err != nil {
			return fmt.Errorf("writing message %d header %s: %w", len(ms)+1, fName, err)
		}
		ct := m.MailHeader.Get("Content-Type")
		mt, mtp, err := mime.ParseMediaType(ct)
		if err != nil {
			return fmt.Errorf("reading message %d header %s content type : %w", len(ms)+1, fName, err)
		}
		switch mt {
		case "multipart/alternative":
			mpw := multipart.NewWriter(f)
			if err := mpw.SetBoundary(mtp["Boundary"]); err != nil {
				return fmt.Errorf("multipart boundary error message %d header %s content type: %w", len(ms)+1, fName, err)
			}
			for _, mb := range m.MailBodies {
				mpwf, err := mpw.CreatePart(mb.Header())
				if err != nil {
					return fmt.Errorf("creating part for multipart message %d body %s: %w", len(ms)+1, fName, err)
				}
				if _, err := io.Copy(mpwf, mb.Reader()); err != nil && !errors.Is(err, io.EOF) {
					return fmt.Errorf("writing message %d body %s: %w", len(ms)+1, fName, err)
				}
			}
		default:
			for _, mb := range m.MailBodies {
				if _, err := io.Copy(f, mb.Reader()); err != nil && !errors.Is(err, io.EOF) {
					return fmt.Errorf("writing message %d body %s: %w", len(ms)+1, fName, err)
				}
			}
		}

	}
	return nil
}
