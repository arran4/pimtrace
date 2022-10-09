package main

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"io"
	"mime/multipart"
	mail2 "net/mail"
	"net/textproto"
	"strconv"
	"time"
)

type MailBody interface {
	Reader() io.Reader
	Header() textproto.MIMEHeader
	FileName() string
	FormName() string
}

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

func (s *MailWithSource) Get(key string) Value {
	// TODO to figure out a good prefix for non-header queries
	return SimpleStringValue(s.MailHeader.Get(key))
}

func (s *MailWithSource) Mail() *MailWithSource {
	return s
}

func (s *MailWithSource) Header() mail.Header {
	return s.MailHeader
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

type SimpleStringValue string

func (s SimpleStringValue) Time() *time.Time {
	t, err := mail2.ParseDate(string(s))
	if err != nil || t.UnixNano() == 0 {
		return nil
	}
	return &t
}

func (s SimpleStringValue) Integer() *int {
	i, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return nil
	}
	ii := int(i)
	return &ii
}

func (s SimpleStringValue) Type() Type {
	return String
}

func (s SimpleStringValue) String() string {
	return string(s)
}

var _ Value = SimpleStringValue("")

type Type int

const (
	String Type = iota
)

type Value interface {
	Type() Type
	String() string
	Time() *time.Time
	Integer() *int
}

type Entry interface {
	Get(string) Value
	Mail() *MailWithSource
	Header() mail.Header
}

type Data interface {
	Len() int
	Entry(n int) Entry
	Mail() []*MailWithSource
}

type PlainOldMailData []*MailWithSource

func (p PlainOldMailData) Len() int {
	return len([]*MailWithSource(p))
}

func (p PlainOldMailData) Entry(n int) Entry {
	if n >= len([]*MailWithSource(p)) || n < 0 {
		return nil
	}
	return ([]*MailWithSource(p))[n]
}

func (p PlainOldMailData) Mail() []*MailWithSource {
	return []*MailWithSource(p)
}

var _ Data = PlainOldMailData(nil)
