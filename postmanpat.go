package postmanpat

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	VERSION          = "0.0.1"
	DEFAULT_BASE_URL = "https://api.postmarkapp.com"
)

type Message struct {
	Subject     string
	HtmlBody    string
	TextBody    string
	To          []*mail.Address
	Attachments []*Attachment
}

func (m *Message) AddAttachment(name string, filename string) error {
	a := &Attachment{}
	a.Name = name
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	a.File = f
	m.Attachments = append(m.Attachments, a)
	return nil
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) AddTo(address string) error {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return err
	}
	m.To = append(m.To, addr)
	return nil
}

type Attachment struct {
	Name string
	File *os.File
}

func joinAddresses(addrs []*mail.Address) string {
	var astr []string
	for _, a := range addrs {
		astr = append(astr, a.String())
	}
	return strings.Join(astr, ",")
}

func (m *Message) Prepare() (*PostmarkMessage, error) {
	pm := &PostmarkMessage{}
	pm.Subject = m.Subject
	pm.To = joinAddresses(m.To)
	pm.HtmlBody = m.HtmlBody
	pm.TextBody = m.TextBody
	for _, a := range m.Attachments {
		pm.AddAttachment(a)
	}
	return pm, nil
}

type PostmarkMessage struct {
	From        string                `json:",omitempty"`
	To          string                `json:",omitempty"`
	Cc          string                `json:",omitempty"`
	Bcc         string                `json:",omitempty"`
	Subject     string                `json:",omitempty"`
	HtmlBody    string                `json:",omitempty"`
	TextBody    string                `json:",omitempty"`
	ReplyTo     string                `json:",omitempty"`
	TrackOpens  bool                  `json:",omitempty"`
	Attachments []*PostmarkAttachment `json:",omitempty"`
}

func (p *PostmarkMessage) AddAttachment(a *Attachment) (bool, error) {
	pa := &PostmarkAttachment{}

	pa.Name = a.Name
	st, err := a.File.Stat()
	if err != nil {
		return false, err
	}
	pa.ContentType = mime.TypeByExtension(filepath.Ext(st.Name()))

	buf := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	defer enc.Close()

	b := bufio.NewReader(a.File)
	if _, err = io.Copy(enc, b); err != nil {
		return false, err
	}

	pa.Content = buf.String()
	p.Attachments = append(p.Attachments, pa)
	return true, nil
}

type PostmarkAttachment struct {
	Name        string
	Content     string
	ContentType string
}

type PostmarkResponse struct {
	To          string
	SubmittedAt time.Time
	MessageID   string
	ErrorCode   int
	Message     string
}

type PostmarkOpenHook struct {
	MessageID string
	Client    struct {
		Name    string
		Company string
		Family  string
	}
	OS struct {
		Name    string
		Company string
		Family  string
	}
	Platform string
}
