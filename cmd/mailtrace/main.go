package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
)

var (
	//inputType = flag.String("input-type", "mbox", "The input type")
	inputFile = flag.String("input", "-", "Input file or - for stdin")
	//outputType = flag.String("output-type", "mbox", "The input type")
	outputFile = flag.String("output", "-", "Output file or - for stdin")
)

type MailBody interface {
	io.Reader
	Header() textproto.MIMEHeader
	FileName() string
	FormName() string
}

type MailWithSource struct {
	MailHeader mail.Header
	MailBodies []*MailBody
	SourceType string
	SourceFile string
}

func main() {
	flag.Parse()
	mails := []*MailWithSource{}
	switch *inputFile {
	case "-":
		nm, err := ReadMbox(os.Stdin, "mbox", *inputFile)
		if err != nil {
			log.Panicln(err)
		}
		mails = append(mails, nm...)
	default:
		nm, err := ReadMboxFile("mbox", *inputFile)
		if err != nil {
			log.Panicln(err)
		}
		mails = append(mails, nm...)
	}
}

func ReadMboxFile(fType, fName string) ([]*MailWithSource, error) {
	f, err := os.OpenFile(fName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("reading Mbox %s: %w", fName, err)
	}
	defer f.Close()
	return ReadMbox(f, fType, fName)
}

func ReadMbox(f io.Reader, fType string, fName string) ([]*MailWithSource, error) {
	ms := []*MailWithSource{}
	for {
		msg, err := mail.ReadMessage(f)
		if err != nil {
			return nil, fmt.Errorf("reading message %d from Mbox %s: %w", len(ms)+1, fName, err)
		}
		if msg == nil {
			return ms, nil
		}
		mb := []*MailBody{}
		ct := msg.Header.Get("Content-Type")
		mt, mtp, err := mime.ParseMediaType(ct)
		switch mt {
		case "multipart/alternative":
			br := multipart.NewReader(msg.Body, mtp["boundary"])
			p, err := br.NextPart()
			if err != nil {
				return nil, fmt.Errorf("reading message %d part %d from Mbox %s: %w", len(ms)+1, len(mb)+1, fName, err)
			}
			mb = append(mb, &MailBody{
				p,
			})
			p.
		default:
			mb = append(mb, &MailBody{
				&multipart.Part{

				},
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
