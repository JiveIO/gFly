package mail

import (
	"app/core/log"
	"app/core/utils"
	"fmt"
	"net/smtp"
)

type Envelop struct {
	To      []string // Required
	ReplyTo []string
	Bcc     []string
	Cc      []string
	Subject string // Required
	Text    string // Required
	HTML    string // Required
}

func Send(envelop Envelop) {
	e := New()
	e.From = fmt.Sprintf("%s <%s>", utils.Getenv("MAIL_NAME", ""), utils.Getenv("MAIL_SENDER", ""))

	if len(envelop.ReplyTo) == 0 {
		e.ReplyTo = []string{utils.Getenv("MAIL_SENDER", "")}
	} else {
		e.ReplyTo = envelop.ReplyTo
	}

	if len(envelop.Bcc) > 0 {
		e.Bcc = envelop.Bcc
	}

	if len(envelop.Cc) > 0 {
		e.Cc = envelop.Cc
	}

	e.To = envelop.To
	e.Subject = envelop.Subject
	e.Text = []byte(envelop.Text)
	e.HTML = []byte(envelop.HTML)

	host := utils.Getenv("MAIL_HOST", "")
	address := fmt.Sprintf("%s:%d", host, utils.Getenv("MAIL_PORT", 587))
	username := utils.Getenv("MAIL_USERNAME", "")
	password := utils.Getenv("MAIL_PASSWORD", "")

	err := e.Send(address, smtp.PlainAuth("", username, password, host))
	if err != nil {
		log.Errorf("Error send mail %v", err)
	}
}
