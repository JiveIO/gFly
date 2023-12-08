package notification

import (
	"app/core/errors"
	"app/core/log"
	"app/core/utils"
	"reflect"
	"sync"
	"time"
)

// -------------------------------------------- Mail Notification ------------------------------------------------------

type MailNotificationData struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

type IMailNotification interface {
	ToEmail() MailNotificationData
}

// --------------------------------------------- DB Notification -------------------------------------------------------

type DatabaseNotificationData struct {
	Type            string
	NotifiableGroup string
	NotifiableId    string
	Data            interface{}
}

type IDatabaseNotification interface {
	ToDatabase() DatabaseNotificationData
}

// -------------------------------------------- SMS Notification -------------------------------------------------------

type SMSNotificationData struct {
	From    string
	To      string
	Message string
}

type ISMSNotification interface {
	ToSMS() SMSNotificationData
}

// ------------------------------------------- Slack Notification ------------------------------------------------------

type SlackNotificationData struct {
	Message string
}

type ISlackNotification interface {
	ToSlack() SlackNotificationData
}

type INotification interface {
	ToEmail() MailNotificationData
	ToDatabase() DatabaseNotificationData
	ToSMS() SMSNotificationData
	ToSlack() SlackNotificationData
}

// ------------------------------------------- Notification Handler ----------------------------------------------------

type INotifiable interface {
	Notify()
}

// --------------------------------------- Notification Orchestration --------------------------------------------------

// Send Deliver many notification types SMS|Mail|Slack|Database.
func Send(notification any) error {
	if !utils.Getenv("NOTIFICATION_ENABLE", false) {
		log.Warnf("[STOP] Notification at %v", time.Now())

		return nil
	}

	var notificationHandlers []INotifiable

	// Check interface type
	mailNotification := reflect.TypeOf((*IMailNotification)(nil)).Elem()
	databaseNotification := reflect.TypeOf((*IDatabaseNotification)(nil)).Elem()
	smsNotification := reflect.TypeOf((*ISMSNotification)(nil)).Elem()
	slackNotification := reflect.TypeOf((*ISlackNotification)(nil)).Elem()

	if reflect.TypeOf(notification).Implements(mailNotification) {
		notificationHandlers = append(notificationHandlers, NewMailNotificationHandler(notification.(IMailNotification)))
	}

	if reflect.TypeOf(notification).Implements(databaseNotification) {
		notificationHandlers = append(notificationHandlers, NewDBNotificationHandler(notification.(IDatabaseNotification)))
	}

	if reflect.TypeOf(notification).Implements(smsNotification) {
		notificationHandlers = append(notificationHandlers, NewSMSNotificationHandler(notification.(ISMSNotification)))
	}

	if reflect.TypeOf(notification).Implements(slackNotification) {
		notificationHandlers = append(notificationHandlers, NewSlackNotificationHandler(notification.(ISlackNotification)))
	}

	if len(notificationHandlers) == 0 {
		return errors.NotYetImplemented
	}

	var wg sync.WaitGroup

	startTime := time.Now()
	// Send message to each channel
	for _, notificationHandler := range notificationHandlers {
		wg.Add(1)
		notificationHandler := notificationHandler
		go func() {
			defer wg.Done()
			notificationHandler.Notify()
		}()
	}

	wg.Wait()

	log.Infof("[RUN] Notification time %v", time.Since(startTime))

	return nil
}
