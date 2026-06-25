package redpanda

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
}

// NewConsumer creates a consumer bound to a consumer group.
// The group is registered on the broker on the first call to Poll.
func NewConsumer(groupID string, topics []string) (*Consumer, error) {
	brokers := common.GetEnvArray("BROKERS", []string{})

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topics...),
		kgo.AutoCommitMarks(), // only commits what is explicitly marked
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("redpanda: new consumer: %w", err)
	}

	return &Consumer{client: client}, nil
}

// Poll starts the consumption loop. For each received message, it calls handler.
// If handler returns nil, the offset is marked for commit.
// If handler returns an error, the message is NOT committed — it will be redelivered.
// Returns when ctx is cancelled (graceful shutdown).
func (c *Consumer) Poll(ctx context.Context, handler func(context.Context, Event) error) error {
	for {
		fetches := c.client.PollFetches(ctx)

		if ctx.Err() != nil {
			return nil
		}

		if fetches.IsClientClosed() {
			return nil
		}

		var fetchErr error
		fetches.EachError(func(topic string, partition int32, err error) {
			fetchErr = fmt.Errorf("redpanda: fetch topic=%s partition=%d: %w", topic, partition, err)
		})
		if fetchErr != nil {
			return fetchErr
		}

		fetches.EachRecord(func(record *kgo.Record) {
			event := eventFromRecord(record)

			if err := handler(ctx, event); err != nil {
				// do not mark commit — message will be redelivered
				return
			}

			c.client.MarkCommitRecords(record)
		})

		// flush marked commits
		commitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.client.CommitMarkedOffsets(commitCtx)
		cancel()
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}

// eventFromRecord reconstructs an Event from a kgo.Record,
// reading fields from the headers and payload from the Value.
func eventFromRecord(record *kgo.Record) Event {
	headers := make(map[string]string, len(record.Headers))
	for _, h := range record.Headers {
		headers[h.Key] = string(h.Value)
	}

	ts, _ := time.Parse(time.RFC3339, headers["ts"])

	var payload struct {
		InventoryID string `json:"inventory_id"`
		ProductID   string `json:"product_id"`
		SKU         string `json:"sku"`
		Quantity    int    `json:"quantity"`
		Warehouse   string `json:"warehouse"`
	}
	json.Unmarshal(record.Value, &payload)

	return Event{
		ID:          string(record.Key),
		TenantID:    headers["tenant_id"],
		Source:      headers["source"],
		Type:        headers["type"],
		Timestamp:   ts,
		InventoryID: payload.InventoryID,
		ProductID:   payload.ProductID,
		SKU:         payload.SKU,
		Quantity:    payload.Quantity,
		Warehouse:   payload.Warehouse,
	}
}
