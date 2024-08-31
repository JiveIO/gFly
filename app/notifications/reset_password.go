package notifications

import (
	notifyMail "github.com/gflydev/notification/mail"
)

type ResetPassword struct {
	Email string
}

func (n ResetPassword) ToEmail() notifyMail.Data {
	return notifyMail.Data{
		To:      n.Email,
		Subject: "gFly - Reset password",
		Body:    "New password was created",
	}
}
