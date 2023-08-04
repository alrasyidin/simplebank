package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAdress   = "smtp.gmail.com"
	smtpServerAdress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	gs := new(GmailSender)
	gs.name = name
	gs.fromEmailAddress = fromEmailAddress
	gs.fromEmailPassword = fromEmailPassword
	return gs
}

func (gm *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", gm.name, gm.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", gm.fromEmailAddress, gm.fromEmailPassword, smtpAuthAdress)
	return e.Send(smtpServerAdress, smtpAuth)
}
