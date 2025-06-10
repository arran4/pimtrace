package maildata

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emersion/go-message/mail"
	"io"
	"mime/multipart"
	"net/textproto"
	"pimtrace"
	"strings"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
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

var _ pimtrace.Entry = (*MailWithSource)(nil)
var _ pimtrace.HasStringArray = (*MailWithSource)(nil)

func (s *MailWithSource) Self() *MailWithSource {
	return s
}

func (s *MailWithSource) HeadersStringArray() (result []string) {
	result = make([]string, 0, s.MailHeader.Len())
	for h := range s.MailHeader.Map() {
		result = append(result, h)
	}
	return
}

func (s *MailWithSource) StringArray(header []string) (result []string) {
	for _, v := range header {
		result = append(result, s.MailHeader.Get(v))
	}
	return
}

func (s *MailWithSource) Get(key string) (pimtrace.Value, error) {
	ks := strings.SplitN(key, ".", 2)
	switch ks[0] {
	case "sz", "sized":
		if len(ks) > 1 {
			v, err := s.Get(ks[1])
			if err != nil {
				return nil, err
			}
			return pimtrace.SimpleIntegerValue(v.Length()), nil
		}
		size := 0
		for k, vs := range s.MailHeader.Map() {
			size += len(k)
			for _, v := range vs {
				size += len(v)
			}
		}
		for _, b := range s.MailBodies {
			switch mb := b.(type) {
			case *MailBodyGeneral:
				size += mb.Body.Len()
			case *MailBodyFromPart:
				if mb.MailBodyGeneral != nil && mb.MailBodyGeneral.Body != nil {
					size += mb.MailBodyGeneral.Body.Len()
				}
			}
		}
		return pimtrace.SimpleIntegerValue(size), nil
	case "h", "header":
		ks = ks[1:]
		fallthrough
	default:
		if len(ks) > 0 {
			return pimtrace.SimpleStringValue(s.MailHeader.Get(ks[0])), nil
		}
		return nil, fmt.Errorf("mail get %w, %s", ErrKeyNotFound, key)
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

type Data []*MailWithSource

func (mdt Data) Truncate(n int) pimtrace.Data {
	mdt = (([]*MailWithSource)(mdt))[:n]
	return mdt
}

func (mdt Data) SetEntry(n int, entry pimtrace.Entry) pimtrace.Data {
	for n > len(mdt) {
		mdt = append((([]*MailWithSource)(mdt)), nil)
	}
	if n == len(mdt) {
		mdt = append(mdt, entry.(*MailWithSource))
	} else {
		(([]*MailWithSource)(mdt))[n] = entry.(*MailWithSource)
	}
	return mdt

}

func (mdt Data) Len() int {
	return len([]*MailWithSource(mdt))
}

func (mdt Data) Entry(n int) pimtrace.Entry {
	if n >= len([]*MailWithSource(mdt)) || n < 0 {
		return nil
	}
	return ([]*MailWithSource(mdt))[n]
}

func (mdt Data) Self() []*MailWithSource {
	return []*MailWithSource(mdt)
}

func (mdt Data) NewSelf() pimtrace.Data {
	return Data(make([]*MailWithSource, 0))
}

var _ pimtrace.Data = Data(nil)

type MailBody interface {
	Reader() io.Reader
	Header() textproto.MIMEHeader
	FileName() string
	FormName() string
}
