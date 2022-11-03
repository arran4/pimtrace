package maildata

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"io"
	"mime/multipart"
	"net/textproto"
	"pimtrace"
	"strings"
	"time"
)

type MailBodyFromPart struct {
	*MailBodyGeneral
	Part *multipart.Part
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
	Body    *bytes.Buffer
	Message *MailWithSource
}

func (m *MailBodyGeneral) Reader() io.Reader {
	return bytes.NewReader(m.Body.Bytes())
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

func (s *MailWithSource) Self() *MailWithSource {
	return s
}

func (s *MailWithSource) Get(key string) pimtrace.Value {
	ks := strings.SplitN(key, ".", 2)
	switch ks[0] {
	//case "sz", "sized": TODO
	//	return SimpleNumberValue(s.
	case "h", "header":
		fallthrough
	default:
		if len(ks) > 1 {
			return pimtrace.SimpleStringValue(s.MailHeader.Get(ks[1]))
		}
		return nil
	}
}

func (s *MailWithSource) Header() *mail.Header {
	return &s.MailHeader
}

func (s *MailWithSource) From() string {
	if f := s.MailHeader.Get("From"); f != "" {
		a, err := mail.ParseAddress(f)
		if err == nil && a != nil {
			return a.Name
		}
	}
	return "nobody"
}

func (s *MailWithSource) Time() time.Time {
	d, _ := s.MailHeader.Date()
	return d
}

type MailDataType []*MailWithSource

func (p MailDataType) Truncate(n int) pimtrace.Data[*MailWithSource] {
	p = (([]*MailWithSource)(p))[:n]
	return p
}

func (p MailDataType) SetEntry(n int, entry pimtrace.Entry[*MailWithSource]) {
	cast := ([]*MailWithSource)(p)
	cast[n] = entry.(*MailWithSource)
}

func (p MailDataType) Len() int {
	return len([]*MailWithSource(p))
}

func (p MailDataType) Entry(n int) pimtrace.Entry[*MailWithSource] {
	if n >= len([]*MailWithSource(p)) || n < 0 {
		return nil
	}
	return ([]*MailWithSource(p))[n]
}

func (p MailDataType) Self() []*MailWithSource {
	return []*MailWithSource(p)
}

var _ pimtrace.Data[*MailWithSource] = MailDataType(nil)

type MailBody interface {
	Reader() io.Reader
	Header() textproto.MIMEHeader
	FileName() string
	FormName() string
}
