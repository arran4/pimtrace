package maildata

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message/mail"
	"github.com/jhillyerd/enmime"
	"io"
	"log"
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
	msg, err := enmime.ReadEnvelope(ff)
	if err != nil {
		return nil, fmt.Errorf("reading message from mail file %s: %w", fName, err)
	}
	if msg == nil || msg.Root.Header == nil {
		return nil, nil
	}
	mws := &MailWithSource{
		SourceFile: fName,
		SourceType: fType,
		MailHeader: mail.HeaderFromMap(msg.Root.Header),
		MailBodies: []MailBody{},
	}
	mws.MailBodies = append(mws.MailBodies, &MailBodyGeneral{
		Body:    bytes.NewBufferString(msg.HTML),
		Message: mws,
	})
	mws.MailBodies = append(mws.MailBodies, &MailBodyGeneral{
		Body:    bytes.NewBufferString(msg.Text),
		Message: mws,
	})

	return []*MailWithSource{mws}, nil

}
