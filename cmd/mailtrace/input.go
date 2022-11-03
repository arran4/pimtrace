package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"pimtrace"
	"pimtrace/dataformats/maildata"
)

func InputHandler(inputType string, inputFile string) (pimtrace.Data, error) {
	mails := []*maildata.MailWithSource{}
	switch inputType {
	case "mailfile":
		switch inputFile {
		case "-":
			nm, err := ReadMailStream(os.Stdin, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := ReadMailFile(inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "mbox":
		switch inputFile {
		case "-":
			nm, err := ReadMBoxStream(os.Stdin, inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := ReadMBoxFile(inputType, inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "list":
		fmt.Println("`--input-type`s: ")
		fmt.Printf(" =%-20s - %s\n", "mailfile", "A single mail file")
		fmt.Printf(" =%-20s - %s\n", "mbox", "Mbox file")
		fmt.Printf(" =%-20s - %s\n", "list", "This help text")
		fmt.Println()
	default:
		fmt.Println("Please specify a -input-type")
		fmt.Println()
	}
	return maildata.Data(mails), nil
}

func ReadMBoxFile(fType, fName string) ([]*maildata.MailWithSource, error) {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return ReadMBoxStream(f, fType, fName)
}

func ReadMBoxStream(f io.Reader, fType string, fName string) ([]*maildata.MailWithSource, error) {
	mbr := mbox.NewReader(f)
	ms := []*maildata.MailWithSource{}
	for {
		mr, err := mbr.NextMessage()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading message %d from Mbox %s: %w", len(ms)+1, fName, err)
		}
		if mr == nil {
			return ms, nil
		}
		mrms, err := ReadMailStream(mr, fType, fName)
		if err != nil {
			return nil, fmt.Errorf("parsing message %d from Mbox %s: %w", len(ms)+1, fName, err)
		}
		ms = append(ms, mrms...)
	}
}

func ReadMailFile(fType, fName string) ([]*maildata.MailWithSource, error) {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading Mbox %s: %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s: %s", fName, err)
		}
	}()
	return ReadMailStream(f, fType, fName)
}

func ReadMailStream(f io.Reader, fType string, fName string) ([]*maildata.MailWithSource, error) {
	ms := []*maildata.MailWithSource{}
	for {
		msg, err := message.Read(f)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading message %d from mail file %s: %w", len(ms)+1, fName, err)
		}
		if msg == nil {
			return ms, nil
		}
		mws := &maildata.MailWithSource{
			SourceFile: fName,
			SourceType: fType,
			MailHeader: mail.HeaderFromMap(msg.Header.Map()),
			MailBodies: []maildata.MailBody{},
		}
		ct := msg.Header.Get("Content-Type")
		mt, mtp, _ := mime.ParseMediaType(ct)
		switch mt {
		case "multipart/alternative":
			br := multipart.NewReader(msg.Body, mtp["boundary"])
			for {
				p, err := br.NextPart()
				if err != nil && !errors.Is(err, io.EOF) {
					return nil, fmt.Errorf("reading message %d part %d %s: %w", len(ms)+1, len(mws.MailBodies)+1, fName, err)
				}
				if p == nil {
					break
				}
				b := bytes.NewBuffer(nil)
				if _, err := io.Copy(b, p); err != nil {
					return nil, fmt.Errorf("reading body of message %d part %d %s: %w", len(ms)+1, len(mws.MailBodies)+1, fName, err)
				}
				mws.MailBodies = append(mws.MailBodies, &maildata.MailBodyFromPart{
					MailBodyGeneral: &maildata.MailBodyGeneral{
						Body:    b,
						Message: mws,
					},
					Part: p,
				})
			}
		default:
			b := bytes.NewBuffer(nil)
			if _, err := io.Copy(b, msg.Body); err != nil {
				return nil, fmt.Errorf("reading body of message %d part %d %s: %w", len(ms)+1, len(mws.MailBodies)+1, fName, err)
			}
			mws.MailBodies = append(mws.MailBodies, &maildata.MailBodyGeneral{
				Body:    b,
				Message: mws,
			})
		}
		ms = append(ms, mws)
	}
}
