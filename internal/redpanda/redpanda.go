package redpanda

import (
	"context"
	"time"

	"github.com/Ricardo-Ronchini/webhook-event-processor/common"
	"github.com/twmb/franz-go/pkg/kgo"
)

const publishTimeout = 10 * time.Second

type Events interface {
	Topic() string
	Key() []byte
	Payload() []byte
	Headers() map[string]string
}

type Redpanda struct {
	producer *Producer
}

func NewRedpanda() *Redpanda {
	brokers := common.GetEnvArray("BROKERS", []string{})

	client, err := NewClient(brokers)
	if err != nil {
		panic(err)
	}

	producer := NewProducer(client)

	return &Redpanda{
		producer: producer,
	}
}

func (r *Redpanda) PublishTopic(ctx context.Context, e Events) error {
	record := &kgo.Record{
		Topic: e.Topic(),
		Key:   e.Key(),
		Value: e.Payload(),
	}

	for k, v := range e.Headers() {
		record.Headers = append(record.Headers, kgo.RecordHeader{
			Key:   k,
			Value: []byte(v),
		})
	}

	publishCtx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	return r.producer.client.ProduceSync(publishCtx, record).FirstErr()
}
