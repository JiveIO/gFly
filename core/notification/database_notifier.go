package notification

import (
	mb "app/core/fluentmodel"
	"app/core/log"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// -------------------------------------- Database Notification Handler ------------------------------------------------

func NewDBNotificationHandler(notification IDatabaseNotification) *DBNotificationHandler {
	return &DBNotificationHandler{
		Notification: notification,
	}
}

type DBNotificationHandler struct {
	Notification IDatabaseNotification
}

func (h *DBNotificationHandler) Notify() {
	data := h.Notification.ToDatabase()
	marshal, err := json.Marshal(data.Data)
	if err != nil {
		log.Errorf("Cause error %v", err)
	}

	db := mb.Instance()

	// Define query string.
	query := `INSERT INTO notifications(id, type, notifiable_group, notifiable_id, data, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Send query to database.
	err = db.Raw(
		query,
		uuid.New(), data.Type, data.NotifiableGroup, data.NotifiableId, marshal, time.Now(), time.Now(),
	).Create(&Notification{})

	if err != nil {
		log.Errorf("Could not create an database notification. Cause error %v", err)
	}

	log.Infof("Send via Database data %v", data)
}
