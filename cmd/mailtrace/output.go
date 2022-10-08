package main

import (
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message/textproto"
	"io"
	"log"
	"os"
)

func OutputHandler(mails []*MailWithSource) error {
	switch *inputType {
	case "mailfile":
		switch *inputFile {
		case "-":
			err := WriteMailStream(mails, os.Stdin, *inputType, *inputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteMailFile(mails, *inputType, *inputFile)
			if err != nil {
				return err
			}
		}
	case "mbox":
		switch *inputFile {
		case "-":
			err := WriteMBoxStream(mails, os.Stdin, *inputType, *inputFile)
			if err != nil {
				return err
			}
		default:
			err := WriteMBoxFile(mails, *inputType, *inputFile)
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

func WriteMBoxFile(ms []*MailWithSource, fType, fName string) error {
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return WriteMBoxStream(ms, f, fType, fName)
}

func WriteMBoxStream(ms []*MailWithSource, f io.Writer, fType string, fName string) error {
	mbw := mbox.NewWriter(f)
	for mi, m := range ms {
		mw, err := mbw.CreateMessage(m.From(), m.Time())
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("creating message %d to Mbox %s: %w", mi+1, fName, err)
		}
		if err := WriteMailStream(ms[mi:mi+1], mw, fType, fName); err != nil {
			return fmt.Errorf("writing message %d to Mbox %s: %w", mi+1, fName, err)
		}
	}
	return nil
}

func WriteMailFile(ms []*MailWithSource, fType, fName string) error {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("reading Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return WriteMailStream(ms, f, fType, fName)
}

func WriteMailStream(ms []*MailWithSource, f io.Writer, fType string, fName string) error {
	for _, m := range ms {
		if err := textproto.WriteHeader(f, textproto.HeaderFromMap(m.MailHeader.Map())); err != nil {
			return
		}
		if err := m.WriteBody(f); err != nil {
			return
		}
	}
	return nil
}
