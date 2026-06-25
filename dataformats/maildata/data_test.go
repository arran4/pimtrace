package maildata

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"github.com/google/go-cmp/cmp"
	"mime/multipart"
	"pimtrace"
	"strings"
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

func TestMailBodyGeneral_Methods(t *testing.T) {
	mb := &MailBodyGeneral{
		Body: bytes.NewBufferString("test body"),
	}

	b := make([]byte, 10)
	n, _ := mb.Reader().Read(b)
	if string(b[:n]) != "test body" {
		t.Errorf("Reader() failed")
	}

	if h := mb.Header(); len(h) != 0 {
		t.Errorf("Header() should be empty")
	}

	if fn := mb.FileName(); fn != "" {
		t.Errorf("FileName() should be empty")
	}

	if formn := mb.FormName(); formn != "" {
		t.Errorf("FormName() should be empty")
	}
}

func TestMailBodyFromPart_Methods(t *testing.T) {
	msg := "Content-Disposition: form-data; name=\"field1\"; filename=\"file.txt\"\r\n" +
		"Content-Type: text/plain\r\n\r\n" +
		"file content\r\n"

	body := "--foo\r\n" + msg + "--foo--\r\n"
	r := multipart.NewReader(bytes.NewReader([]byte(body)), "foo")
	part, err := r.NextPart()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	mb := &MailBodyFromPart{
		MailBodyGeneral: &MailBodyGeneral{},
		Part: part,
	}

	if h := mb.Header(); h.Get("Content-Type") != "text/plain" {
		t.Errorf("Header() failed")
	}

	if fn := mb.FileName(); fn != "file.txt" {
		t.Errorf("FileName() failed")
	}

	if formn := mb.FormName(); formn != "field1" {
		t.Errorf("FormName() failed")
	}
}

func TestMailWithSource_Self(t *testing.T) {
	m := &MailWithSource{}
	if m.Self() != m {
		t.Errorf("Self() didn't return same pointer")
	}
}

func TestMailWithSource_Get(t *testing.T) {
	h := mail.Header{}
	h.Set("Subject", "Hello")
	m := &MailWithSource{
		MailHeader: h,
		MailBodies: []MailBody{nil, nil},
	}

	v, err := m.Get("sz")
	if err != nil {
		t.Errorf("Get(sz) error = %v", err)
	}
	if iv, ok := v.(pimtrace.SimpleIntegerValue); !ok || int(iv) != 2 {
		t.Errorf("Get(sz) expected 2, got %v", v)
	}

	v, err = m.Get("h.Subject")
	if err != nil {
		t.Errorf("Get(h.Subject) error = %v", err)
	}
	if sv, ok := v.(pimtrace.SimpleStringValue); !ok || string(sv) != "Hello" {
		t.Errorf("Get(h.Subject) expected Hello, got %v", v)
	}
}

func TestData(t *testing.T) {
	var d Data = make([]*MailWithSource, 0)

	if d.Len() != 0 {
		t.Errorf("Len() = %v, want 0", d.Len())
	}

	r1 := &MailWithSource{}
	r2 := &MailWithSource{}

	d = d.SetEntry(0, r1).(Data)
	if d.Len() != 1 {
		t.Errorf("Len() = %v, want 1", d.Len())
	}

	d = d.SetEntry(2, r2).(Data) // Should pad
	if d.Len() != 3 {
		t.Errorf("Len() = %v, want 3", d.Len())
	}
	if d.Entry(2) != r2 {
		t.Errorf("Entry(2) != r2")
	}
	if d.Entry(1) != (*MailWithSource)(nil) {
		t.Errorf("Entry(1) should be nil padded")
	}
	if d.Entry(5) != nil {
		t.Errorf("Entry(5) should be nil (out of bounds)")
	}

	d = d.Truncate(2).(Data)
	if d.Len() != 2 {
		t.Errorf("Truncate(2) Len = %v, want 2", d.Len())
	}

	if d.Self() == nil {
		t.Errorf("Self() returned nil")
	}

	dNew := d.NewSelf()
	if dNew.Len() != 0 {
		t.Errorf("NewSelf() Len = %v, want 0", dNew.Len())
	}
}

func TestData_Output(t *testing.T) {
	var d Data = make([]*MailWithSource, 0)
	h := mail.Header{}
	h.Set("Subject", "Test Subj")
	h.Set("From", "tester@example.com")
	h.Set("Date", "Thu, 13 Feb 1969 23:32:54 -0330")

	_ = append(d, &MailWithSource{
		MailHeader: h,
		MailBodies: []MailBody{
			&MailBodyGeneral{
				Body: bytes.NewBufferString("test body content"),
			},
		},
	})

	// CSV and Table format streams wrapper
	//d.WriteCSVFile("-")
	//d.WriteTableFile("-")

	// Mbox and Mail wrapper
	//d.WriteMBoxFile("-")
	//d.WriteMailFile("-")
}

func TestReadMailStream(t *testing.T) {
	mailData := `From: "John Doe" <john@example.com>
To: "Jane Doe" <jane@example.com>
Subject: Test Email
Date: Thu, 13 Feb 1969 23:32:54 -0330
Content-Type: text/plain; charset="utf-8"

This is a test email body.
`

	r := strings.NewReader(mailData)
	res, err := ReadMailStream(r, "mail", "test.eml")
	if err != nil {
		t.Errorf("ReadMailStream error = %v", err)
	}
	if len(res) != 1 {
		t.Errorf("ReadMailStream expected 1 message, got %d", len(res))
	}

	if res[0].SourceType != "mail" || res[0].SourceFile != "test.eml" {
		t.Errorf("ReadMailStream source info wrong")
	}

	if res[0].MailHeader.Get("Subject") != "Test Email" {
		t.Errorf("ReadMailStream subject = %v, want Test Email", res[0].MailHeader.Get("Subject"))
	}
}

func TestReadMBoxStream(t *testing.T) {
	mboxData := `From MAILER-DAEMON Thu Feb 13 23:32:54 1969
From: "John Doe" <john@example.com>
To: "Jane Doe" <jane@example.com>
Subject: Message 1
Date: Thu, 13 Feb 1969 23:32:54 -0330
Content-Type: text/plain; charset="utf-8"

Body 1
From MAILER-DAEMON Thu Feb 13 23:33:54 1969
From: "Alice" <alice@example.com>
To: "Bob" <bob@example.com>
Subject: Message 2
Date: Thu, 13 Feb 1969 23:33:54 -0330
Content-Type: text/plain; charset="utf-8"

Body 2
`

	r := strings.NewReader(mboxData)
	res, err := ReadMBoxStream(r, "mbox", "test.mbox")
	if err != nil {
		t.Errorf("ReadMBoxStream error = %v", err)
	}
	if len(res) != 2 {
		t.Errorf("ReadMBoxStream expected 2 messages, got %d", len(res))
	}
	if res[0].MailHeader.Get("Subject") != "Message 1" {
		t.Errorf("ReadMBoxStream msg1 subject = %v", res[0].MailHeader.Get("Subject"))
	}
	if res[1].MailHeader.Get("Subject") != "Message 2" {
		t.Errorf("ReadMBoxStream msg2 subject = %v", res[1].MailHeader.Get("Subject"))
	}
}

func TestMBoxOutput_Execute(t *testing.T) {
	mo := &MBoxOutput{}
	var d Data = make([]*MailWithSource, 0)

	// Valid data
	res, err := mo.Execute(d, nil)
	if err != nil {
		t.Errorf("MBoxOutput.Execute error = %v", err)
	}
	if res == nil {
		t.Errorf("MBoxOutput.Execute expected original data back")
	}

	// Invalid data
	_, err = mo.Execute(nil, nil)
	if err == nil {
		t.Errorf("MBoxOutput.Execute with nil expected error")
	}
}

func TestMailWithSource_Header(t *testing.T) {
	h := mail.Header{}
	h.Set("Subject", "Test")
	m := &MailWithSource{
		MailHeader: h,
	}
	resH := m.Header()
	if resH.Get("Subject") != "Test" {
		t.Errorf("Header() failed")
	}
}
