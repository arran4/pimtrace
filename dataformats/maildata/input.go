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

func ReadMBoxStream(f io.Reader, fType string, fName string, ops ...any) (res []*MailWithSource, err error) {
	ff, closers, err := dataformats.ReaderStreamMapperOptionProcessor(f, ops)
	defer func() {
		for i := range closers {
			fc := closers[len(closers)-i-1]
			if cerr := fc.Close(); cerr != nil {
				if err == nil {
					err = fmt.Errorf("closing ReaderStreamMapper: %w", cerr)
				} else {
					log.Printf("error closing ReaderStreamMapper: %s", cerr)
				}
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	mbr := mbox.NewReader(ff)
	var ms []*MailWithSource
	for {
		mr, nextErr := mbr.NextMessage()
		if nextErr != nil && !errors.Is(nextErr, io.EOF) {
			return nil, fmt.Errorf("reading message %d from Mbox %s: %w", len(ms)+1, fName, nextErr)
		}
		if mr == nil {
			return ms, nil
		}
		mrms, readErr := ReadMailStream(mr, fType, fName)
		if readErr != nil {
			log.Printf("parsing message %d from Mbox %s: %v", len(ms)+1, fName, readErr)
			continue
		}
		ms = append(ms, mrms...)
	}
}

func ReadMailStream(f io.Reader, fType string, fName string, ops ...any) (res []*MailWithSource, err error) {
	ff, closers, err := dataformats.ReaderStreamMapperOptionProcessor(f, ops)
	defer func() {
		for i := range closers {
			fc := closers[len(closers)-i-1]
			if cerr := fc.Close(); cerr != nil {
				if err == nil {
					err = fmt.Errorf("closing ReaderStreamMapper: %w", cerr)
				} else {
					log.Printf("error closing ReaderStreamMapper: %s", cerr)
				}
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	msg, err := enmime.ReadEnvelope(ff)
	if err != nil {
		return nil, fmt.Errorf("reading message: %w", err)
	}
	if msg == nil || msg.Root.Header == nil {
		return nil, nil
	}
	mws := &MailWithSource{
		SourceFile: fName,
		SourceType: fType,
		MailHeader: mail.HeaderFromMap(msg.Root.Header),
	}
	mws.MailBodies = []MailBody{
		&MailBodyGeneral{
			Body:    bytes.NewBufferString(msg.HTML),
			Message: mws,
		}, &MailBodyGeneral{
			Body:    bytes.NewBufferString(msg.Text),
			Message: mws,
		},
	}

	return []*MailWithSource{mws}, nil

}
