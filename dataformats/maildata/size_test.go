package maildata

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"testing"
)

func TestMailGetSize(t *testing.T) {
	mh := mail.Header{}
	mh.Set("Subject", "Test")
	m := &MailWithSource{
		MailHeader: mh,
		MailBodies: []MailBody{
			&MailBodyGeneral{Body: bytes.NewBufferString("Hello"), Message: nil},
		},
	}
	v, err := m.Get("sz")
	if err != nil {
		t.Fatalf("sz err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 16 {
		t.Errorf("sz expected 16 got %v", v)
	}
	v, err = m.Get("sized.h.Subject")
	if err != nil {
		t.Fatalf("sized err: %v", err)
	}
	if i := v.Integer(); i == nil || *i != 4 {
		t.Errorf("sized expected 4 got %v", v)
	}
}
