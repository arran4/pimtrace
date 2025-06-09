package maildata

import (
	"bytes"
	"testing"
)

func TestMailWithSourceGetSize(t *testing.T) {
	m := &MailWithSource{
		MailBodies: []MailBody{
			&MailBodyGeneral{Body: bytes.NewBufferString("a")},
			&MailBodyGeneral{Body: bytes.NewBufferString("b")},
		},
	}
	v, err := m.Get("sz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i := v.Integer()
	if i == nil || *i != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}
