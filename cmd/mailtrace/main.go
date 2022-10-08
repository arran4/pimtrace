package main

import (
	"bytes"
	"flag"
	"io"
	"log"
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
	mails, err := InputHandler()
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}

	if err := OutputHandler(mails); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}
}
