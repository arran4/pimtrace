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

func OutputHandler(d Data, outputType *string, outputFile *string) error {
	switch *outputType {
	case "mailfile":
		switch *outputFile {
		case "-":
			err := WriteMailStream(d, os.Stdin, *outputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteMailFile(d, *outputFile)
			if err != nil {
				return err
			}
		}
	case "mbox":
		switch *outputFile {
		case "-":
			err := WriteMBoxStream(d, os.Stdin, *outputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteMBoxFile(d, *outputFile)
			if err != nil {
				return err
			}
		}
	case "csv":
		switch *outputFile {
		case "-":
			err := WriteCSVStream(d, os.Stdin, *outputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteCSVFile(d, *outputFile)
			if err != nil {
				return err
			}
		}
	case "count":
		fmt.Println(d.Len())
	case "list":
		fmt.Println("`--output-type`s: ")
		fmt.Printf(" =%-20s - %s\n", "mailfile", "A single mail file")
		fmt.Printf(" =%-20s - %s\n", "mbox", "Mbox file")
		fmt.Printf(" =%-20s - %s\n", "list", "This help text")
		fmt.Printf(" =%-20s - %s\n", "count", "Just a count")
		fmt.Printf(" =%-20s - %s\n", "csv", "Data in csv format")
		fmt.Println()
	default:
		fmt.Println("Please specify a -input-type")
		fmt.Println()
	}
	return nil
}

func WriteMBoxFile(d Data, fName string) error {
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return WriteMBoxStream(d, f, fName)
}

func WriteMBoxStream(d Data, f io.Writer, fName string) error {
	ms := d.Mail()
	mbw := mbox.NewWriter(f)
	for mi, m := range ms {
		mw, err := mbw.CreateMessage(m.From(), m.Time())
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("creating message %d to Mbox %s: %w", mi+1, fName, err)
		}
		if err := WriteMailStream(MailDataType(ms[mi:mi+1]), mw, fName); err != nil {
			return fmt.Errorf("writing message %d to Mbox %s: %w", mi+1, fName, err)
		}
	}
	return nil
}

func WriteMailFile(d Data, fName string) error {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("writing mail file %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return WriteMailStream(d, f, fName)
}

func WriteMailStream(d Data, f io.Writer, fName string) error {
	for mi, m := range d.Mail() {
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
