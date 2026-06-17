package redpanda

import (
	"context"
	"fmt"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
}

// NewConsumer cria um consumer vinculado a um consumer group.
// O group é registrado no broker na primeira chamada a Poll.
func NewConsumer(groupID string, topics []string) (*Consumer, error) {
	brokers := common.GetEnvArray("BROKERS", []string{})

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topics...),
		kgo.AutoCommitMarks(), // só commita o que for explicitamente marcado
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("redpanda: new consumer: %w", err)
	}

	return &Consumer{client: client}, nil
}

// Poll inicia o loop de consumo. Para cada mensagem recebida, chama handler.
// Se handler retornar nil, o offset é marcado para commit.
// Se handler retornar erro, a mensagem NÃO é commitada — será re-entregue.
// Retorna quando ctx for cancelado (graceful shutdown).
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
				// não marca commit — mensagem será re-entregue
				return
			}

			c.client.MarkCommitRecords(record)
		})

		// flush dos commits marcados
		commitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.client.CommitMarkedOffsets(commitCtx)
		cancel()
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}

// eventFromRecord reconstrói um Event a partir de um kgo.Record,
// lendo os campos dos headers e o payload do Value.
func eventFromRecord(record *kgo.Record) Event {
	headers := make(map[string]string, len(record.Headers))
	for _, h := range record.Headers {
		headers[h.Key] = string(h.Value)
	}

	ts, _ := time.Parse(time.RFC3339, headers["ts"])

	return Event{
		ID:        string(record.Key),
		TenantID:  headers["tenant_id"],
		Source:    headers["source"],
		Type:      headers["type"],
		Timestamp: ts,
		Data:      record.Value,
	}
}
