package redpanda

import (
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
)

type Event struct {
	ID        string
	TenantID  string
	Source    string
	Type      string
	Timestamp time.Time
	Data      []byte
}

func (e Event) Topic() string {
	return common.MessengerTopicID
}

func (e Event) Key() []byte {
	return []byte(e.ID)
}

func (e Event) Payload() []byte {
	return e.Data
}

func (e Event) Headers() map[string]string {
	return map[string]string{
		"source":    e.Source,
		"type":      e.Type,
		"tenant_id": e.TenantID,
		"ts":        e.Timestamp.Format(time.RFC3339),
	}
}
