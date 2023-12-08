package notification

import (
	"app/core/log"
	"app/core/utils"
	"encoding/json"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// --------------------------------------- SMS Notification Handler ---------------------------------------------------

func NewSMSNotificationHandler(notification ISMSNotification) *SMSNotificationHandler {
	return &SMSNotificationHandler{
		Notification: notification,
	}
}

type SMSNotificationHandler struct {
	Notification ISMSNotification
}

// Notify Send SMS message.
//
// Note: Need to verify.
//
// Refer:
//
//	https://www.twilio.com/blog/send-sms-30-seconds-golang
//	https://www.twilio.com/docs/messaging/quickstart/go
//	https://github.com/twilio/twilio-go
func (h *SMSNotificationHandler) Notify() {
	data := h.Notification.ToSMS()

	accountSid := utils.Getenv("TWILIO_ACCOUNT_SID", "")
	authToken := utils.Getenv("TWILIO_AUTH_TOKEN", "")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &api.CreateMessageParams{}
	params.SetTo(data.To)
	params.SetFrom(data.From)
	params.SetBody(data.Message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Errorf("Error sending SMS message %v. Data %v", err, data)
	}

	response, _ := json.Marshal(*resp)

	log.Infof("Not yet implemented! Send via SMS data %v. Response %s", data, utils.UnsafeString(response))
}
