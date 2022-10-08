package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
)

func InputHandler() ([]*MailWithSource, error) {
	mails := []*MailWithSource{}
	switch *inputType {
	case "mailfile":
		switch *inputFile {
		case "-":
			nm, err := ReadMailStream(os.Stdin, *inputType, *inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := ReadMailFile(*inputType, *inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		}
	case "mbox":
		switch *inputFile {
		case "-":
			nm, err := ReadMBoxStream(os.Stdin, *inputType, *inputFile)
			if err != nil {
				return nil, err
			}
			mails = append(mails, nm...)
		default:
			nm, err := ReadMBoxFile(*inputType, *inputFile)
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
	return mails, nil
}

func ReadMBoxFile(fType, fName string) ([]*MailWithSource, error) {
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

func ReadMBoxStream(f io.Reader, fType string, fName string) ([]*MailWithSource, error) {
	mbr := mbox.NewReader(f)
	ms := []*MailWithSource{}
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

func ReadMailFile(fType, fName string) ([]*MailWithSource, error) {
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

func ReadMailStream(f io.Reader, fType string, fName string) ([]*MailWithSource, error) {
	ms := []*MailWithSource{}
	for {
		msg, err := mail.ReadMessage(f)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading message %d from mail file %s: %w", len(ms)+1, fName, err)
		}
		if msg == nil {
			return ms, nil
		}
		mb := []MailBody{}
		ct := msg.Header.Get("Content-Type")
		mt, mtp, err := mime.ParseMediaType(ct)
		switch mt {
		case "multipart/alternative":
			br := multipart.NewReader(msg.Body, mtp["boundary"])
			for {
				p, err := br.NextPart()
				if err != nil && !errors.Is(err, io.EOF) {
					return nil, fmt.Errorf("reading message %d part %d from Mbox %s: %w", len(ms)+1, len(mb)+1, fName, err)
				}
				if p == nil {
					break
				}
				b := bytes.NewBuffer(nil)
				if _, err := io.Copy(b, p); err != nil {
					return nil, fmt.Errorf("reading body of message %d part %d from Mbox %s: %w", len(ms)+1, len(mb)+1, fName, err)
				}
				mb = append(mb, &MailBodyFromPart{
					Body: b,
					Part: p,
				})
			}
		default:
			b := bytes.NewBuffer(nil)
			if _, err := io.Copy(b, msg.Body); err != nil {
				return nil, fmt.Errorf("reading body of message %d part %d from Mbox %s: %w", len(ms)+1, len(mb)+1, fName, err)
			}
			mb = append(mb, &MailBodyGeneral{
				Body: b,
			})
		}
		ms = append(ms, &MailWithSource{
			SourceFile: fName,
			SourceType: fType,
			MailHeader: msg.Header,
			MailBodies: mb,
		})
	}
}