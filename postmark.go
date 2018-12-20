package postmark

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

// Message is an email message
type Message struct {
	Subject     string
	HTMLBody    string
	TextBody    string
	To          []*mail.Address
	Attachments []*Attachment
}

// AddAttachment is a helper to add a file attachment to an email
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

// NewMessage creates a new mail message
func NewMessage() *Message {
	return &Message{}
}

// AddTo adds a new recipient to the message
func (m *Message) AddTo(address string) error {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return err
	}
	m.To = append(m.To, addr)
	return nil
}

// Attachment is a Message attachnment backed by a file
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

func newPostmarkMessage(m *Message) (*postmarkMessage, error) {
	pm := &postmarkMessage{}
	pm.Subject = m.Subject
	pm.To = joinAddresses(m.To)
	pm.HTMLBody = m.HTMLBody
	pm.TextBody = m.TextBody
	for _, a := range m.Attachments {
		pm.addAttachment(a)
	}
	return pm, nil
}

type postmarkMessage struct {
	From        string                `json:",omitempty"`
	To          string                `json:",omitempty"`
	Cc          string                `json:",omitempty"`
	Bcc         string                `json:",omitempty"`
	Subject     string                `json:",omitempty"`
	HTMLBody    string                `json:",omitempty"`
	TextBody    string                `json:",omitempty"`
	ReplyTo     string                `json:",omitempty"`
	TrackOpens  bool                  `json:",omitempty"`
	Attachments []*postmarkAttachment `json:",omitempty"`
}

func (p *postmarkMessage) addAttachment(a *Attachment) (bool, error) {
	pa := &postmarkAttachment{}

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

type postmarkAttachment struct {
	Name        string
	Content     string
	ContentType string
}

// Response is recieved by postmark when an email is submitted for sending
type Response struct {
	To          string
	SubmittedAt time.Time
	MessageID   string
	ErrorCode   int
	Message     string
}
