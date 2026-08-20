package workflow

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"strings"
)

type Notification struct {
	ID        string
	RecordID  string
	Recipient string
	Channel   string
	Subject   string
	Body      string
	State     string
}

func BuildNotification(record model.Record, recipient, channel string) Notification {
	subject := "claim update: " + record.ID
	body := fmt.Sprintf("batch %s is %s", record.BatchID, record.Status)
	return Notification{ID: record.ID + "-notice", RecordID: record.ID, Recipient: recipient, Channel: channel, Subject: subject, Body: body, State: "queued"}
}

func ValidateNotification(item Notification) error {
	if strings.TrimSpace(item.Recipient) == "" {
		return fmt.Errorf("recipient required")
	}
	switch item.Channel {
	case "console", "file", "webhook":
		return nil
	default:
		return fmt.Errorf("unsupported channel %s", item.Channel)
	}
}

func DeliverNotification(item Notification) (Notification, error) {
	if err := ValidateNotification(item); err != nil {
		return item, err
	}
	if item.State != "queued" {
		return item, fmt.Errorf("notification is %s", item.State)
	}
	item.State = "delivered"
	return item, nil
}

func FailNotification(item Notification, reason string) Notification {
	item.State = "failed:" + reason
	return item
}

func QueueFor(records []model.Record, recipient string) []Notification {
	result := make([]Notification, 0, len(records))
	for _, record := range records {
		if record.Status == "approved" || record.Status == "rejected" || record.Status == "archived" {
			result = append(result, BuildNotification(record, recipient, "console"))
		}
	}
	return result
}
