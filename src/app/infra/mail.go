package infra

import "log"

// NewFakeMail instances a new FakeMail
func NewFakeMail() *FakeMail {
	return &FakeMail{}
}

// FakeMail is fake mail driver. It logs mails to stdout
type FakeMail struct{}

// Send logs mail to stdout and returns nil
func (fm *FakeMail) Send(to []string, subject string, body []byte) error {
	log.Printf("Mail sent! to: %v, subject: %s, body: %s", to, subject, body)
	return nil
}
