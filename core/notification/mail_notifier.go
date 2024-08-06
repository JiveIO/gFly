package notification

import (
	"app/core/log"
	"app/core/mail"
)

// --------------------------------------- Mail Notification Handler ---------------------------------------------------

func NewMailNotificationHandler(notification IMailNotification) *MailNotificationHandler {
	return &MailNotificationHandler{
		Notification: notification,
	}
}

type MailNotificationHandler struct {
	Notification IMailNotification
}

func (h *MailNotificationHandler) Notify() {
	data := h.Notification.ToEmail()

	envelop := mail.Envelop{
		To:      []string{data.To},
		Subject: data.Subject,
		HTML:    data.Body,
	}

	if data.Cc != "" {
		envelop.Cc = []string{data.Cc}
	}

	log.Tracef("Send via Mail data %v", data)

	mail.Send(envelop)
}
