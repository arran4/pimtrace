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
	"pimtrace/dataformats"
)

func ReadMBoxStream(f io.Reader, fType string, fName string, ops ...any) ([]*MailWithSource, error) {
	ff, closers, err := dataformats.ReaderStreamMapperOptionProcessor(f, ops)
	defer func() {
		for _, fc := range closers {
			if err := fc.Close(); err != nil {
				log.Printf("error closing ReaderStreamMapper: %s", err)
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	mbr := mbox.NewReader(ff)
	var ms []*MailWithSource
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

func ReadMailStream(f io.Reader, fType string, fName string, ops ...any) ([]*MailWithSource, error) {
	var ms []*MailWithSource
	ff, closers, err := dataformats.ReaderStreamMapperOptionProcessor(f, ops)
	defer func() {
		for _, fc := range closers {
			if err := fc.Close(); err != nil {
				log.Printf("error closing ReaderStreamMapper: %s", err)
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	for {
		msg, err := message.Read(ff)
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
