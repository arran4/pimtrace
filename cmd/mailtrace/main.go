package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/emersion/go-mbox"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
)

var (
	inputType = flag.String("input-type", "list", "The input type")
	inputFile = flag.String("input", "-", "Input file or - for stdin")
	//outputType = flag.String("output-type", "list", "The input type")
	outputFile = flag.String("output", "-", "Output file or - for stdin")
)

type MailBody interface {
	io.Reader
	Header() textproto.MIMEHeader
	FileName() string
	FormName() string
}

type MailBodyFromPart struct {
	Body *bytes.Buffer
	Part *multipart.Part
}

func (m *MailBodyFromPart) Read(p []byte) (n int, err error) {
	return m.Body.Read(p)
}

func (m *MailBodyFromPart) Header() textproto.MIMEHeader {
	return m.Part.Header
}

func (m *MailBodyFromPart) FileName() string {
	return m.Part.FileName()
}

func (m *MailBodyFromPart) FormName() string {
	return m.Part.FormName()
}

var _ MailBody = (*MailBodyFromPart)(nil)

type MailBodyGeneral struct {
	Body *bytes.Buffer
}

func (m *MailBodyGeneral) Read(p []byte) (n int, err error) {
	return m.Body.Read(p)
}

func (m *MailBodyGeneral) Header() textproto.MIMEHeader {
	return map[string][]string{}
}

func (m *MailBodyGeneral) FileName() string {
	return ""
}

func (m *MailBodyGeneral) FormName() string {
	return ""
}

var _ MailBody = (*MailBodyGeneral)(nil)

type MailWithSource struct {
	MailHeader mail.Header
	MailBodies []MailBody
	SourceType string
	SourceFile string
}

func main() {
	flag.Parse()
	mails := []*MailWithSource{}
	switch *inputType {
	case "mailfile":
		switch *inputFile {
		case "-":
			nm, err := ReadMailStream(os.Stdin, *inputType, *inputFile)
			if err != nil {
				log.Panicln(err)
			}
			mails = append(mails, nm...)
		default:
			nm, err := ReadMailFile(*inputType, *inputFile)
			if err != nil {
				log.Panicln(err)
			}
			mails = append(mails, nm...)
		}
	case "mbox":
		switch *inputFile {
		case "-":
			nm, err := ReadMBoxStream(os.Stdin, *inputType, *inputFile)
			if err != nil {
				log.Panicln(err)
			}
			mails = append(mails, nm...)
		default:
			nm, err := ReadMBoxFile(*inputType, *inputFile)
			if err != nil {
				log.Panicln(err)
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
}

func ReadMBoxFile(fType, fName string) ([]*MailWithSource, error) {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading Mbox %s: %w", fName, err)
	}
	defer f.Close()
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
	defer f.Close()
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
