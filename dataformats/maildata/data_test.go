package maildata

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
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

func TestMailWithSource_Get_Header(t *testing.T) {
	h := mail.Header{}
	h.Set("Subject", "Test Subject")
	m := &MailWithSource{
		MailHeader: h,
	}

	tests := []struct {
		key      string
		expected string
		err      bool
	}{
		{"h.Subject", "Test Subject", false},
		{"header.Subject", "Test Subject", false},
		{"Subject", "Test Subject", false},
		{"Unknown", "", false}, // header.Get("Unknown") returns "", which is valid result, no error.
	}

	for _, test := range tests {
		v, err := m.Get(test.key)
		if (err != nil) != test.err {
			t.Errorf("Get(%q) error = %v, wantErr %v", test.key, err, test.err)
			continue
		}
		if !test.err {
			if s := v.String(); s != test.expected {
				t.Errorf("Get(%q) = %q, want %q", test.key, s, test.expected)
			}
		}
	}
}

func TestMailWithSource_From(t *testing.T) {
	h := mail.Header{}
	h.Set("From", "Test User <test@example.com>")
	m := &MailWithSource{
		MailHeader: h,
	}

	if got := m.From(); got != "Test User" {
		t.Errorf("From() = %q, want %q", got, "Test User")
	}

	h2 := mail.Header{}
	m2 := &MailWithSource{MailHeader: h2}
	if got := m2.From(); got != "nobody" {
		t.Errorf("From() = %q, want %q", got, "nobody")
	}
}

func TestMailWithSource_Time(t *testing.T) {
	h := mail.Header{}
	now := time.Now().Truncate(time.Second) // mail date resolution is second usually
	h.SetDate(now)
	m := &MailWithSource{
		MailHeader: h,
	}

	if got := m.Time(); !got.Equal(now) {
		t.Errorf("Time() = %v, want %v", got, now)
	}
}

func TestData_Methods(t *testing.T) {
	m1 := &MailWithSource{SourceFile: "1"}
	m2 := &MailWithSource{SourceFile: "2"}
	d := Data{m1, m2}

	if d.Len() != 2 {
		t.Errorf("Len() = %d, want 2", d.Len())
	}

	if e := d.Entry(0); e != m1 {
		t.Errorf("Entry(0) != m1")
	}
	if e := d.Entry(2); e != nil {
		t.Errorf("Entry(2) should be nil")
	}

	d2 := d.Truncate(1)
	if d2.Len() != 1 {
		t.Errorf("Truncate(1) len = %d, want 1", d2.Len())
	}

	d3 := d.NewSelf()
	if d3.Len() != 0 {
		t.Errorf("NewSelf() len = %d, want 0", d3.Len())
	}

	d4 := d.SetEntry(0, m2)
	if d4.Entry(0) != m2 {
		t.Errorf("SetEntry(0) failed")
	}

	d5 := d.SetEntry(3, m1)
	if d5.Len() != 4 {
		t.Errorf("SetEntry(3) len = %d, want 4", d5.Len())
	}
}

func TestMailWithSource_HeadersStringArray(t *testing.T) {
	h := mail.Header{}
	h.Set("A", "1")
	h.Set("B", "2")
	m := &MailWithSource{MailHeader: h}

	got := m.HeadersStringArray()
	// Map iteration order is random, so check existence
	if len(got) != 2 {
		t.Errorf("HeadersStringArray len = %d, want 2", len(got))
	}
}

func TestMailWithSource_StringArray(t *testing.T) {
	h := mail.Header{}
	h.Set("A", "1")
	h.Set("B", "2")
	m := &MailWithSource{MailHeader: h}

	got := m.StringArray([]string{"A", "B", "C"})
	expected := []string{"1", "2", ""}
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Errorf("StringArray mismatch (-want +got):\n%s", diff)
	}
}
