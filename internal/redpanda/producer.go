package redpanda

import "github.com/twmb/franz-go/pkg/kgo"

type Producer struct {
	client *kgo.Client
}

func NewProducer(client *Client) *Producer {
	return &Producer{client: client.Client}
}
