package redpanda

import "github.com/twmb/franz-go/pkg/kgo"

type Client struct {
	*kgo.Client
}

func NewClient(brokers []string) (*Client, error) {
	c, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{c}, nil
}
