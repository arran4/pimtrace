package maildata

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
)

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
		msg, err := message.Read(f)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading message %d from mail file %s: %w", len(ms)+1, fName, err)
		}
		if msg == nil {
			return ms, nil
		}
		mws := &MailWithSource{
			SourceFile: fName,
			SourceType: fType,
			MailHeader: mail.HeaderFromMap(msg.Header.Map()),
			MailBodies: []MailBody{},
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
				mws.MailBodies = append(mws.MailBodies, &MailBodyFromPart{
					MailBodyGeneral: &MailBodyGeneral{
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
			mws.MailBodies = append(mws.MailBodies, &MailBodyGeneral{
				Body:    b,
				Message: mws,
			})
		}
		ms = append(ms, mws)
	}
}
