package notification

import (
	"app/core/log"
	"app/core/notification/slack"
	"app/core/utils"
)

// --------------------------------------- Slack Notification Handler ---------------------------------------------------

func NewSlackNotificationHandler(notification ISlackNotification) *SlackNotificationHandler {
	return &SlackNotificationHandler{
		Notification: notification,
	}
}

type SlackNotificationHandler struct {
	Notification ISlackNotification
}

func (h *SlackNotificationHandler) Notify() {
	data := h.Notification.ToSlack()

	slackConn := slack.New(utils.Getenv("SLACK_CHANNEL", ""))

	// Send notification
	err := slackConn.Notify(data.Message)
	if err != nil {
		log.Errorf("Error %v", err)
	}

	log.Infof("Send via Slack data %v", data)
}
