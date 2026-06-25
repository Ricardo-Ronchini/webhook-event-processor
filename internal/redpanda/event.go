package redpanda

import (
	"encoding/json"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
)

type Event struct {
	ID          string
	TenantID    string
	Source      string
	Type        string
	Timestamp   time.Time
	InventoryID string
	ProductID   string
	SKU         string
	Quantity    int
	Warehouse   string
}

func (e Event) Topic() string {
	return common.MessengerTopicID
}

func (e Event) Key() []byte {
	return []byte(e.ID)
}

func (e Event) Payload() []byte {
	data, _ := json.Marshal(struct {
		InventoryID string `json:"inventory_id"`
		ProductID   string `json:"product_id"`
		SKU         string `json:"sku"`
		Quantity    int    `json:"quantity"`
		Warehouse   string `json:"warehouse,omitempty"`
	}{
		InventoryID: e.InventoryID,
		ProductID:   e.ProductID,
		SKU:         e.SKU,
		Quantity:    e.Quantity,
		Warehouse:   e.Warehouse,
	})
	return data
}

func (e Event) Headers() map[string]string {
	return map[string]string{
		"source":    e.Source,
		"type":      e.Type,
		"tenant_id": e.TenantID,
		"ts":        e.Timestamp.Format(time.RFC3339),
	}
}
