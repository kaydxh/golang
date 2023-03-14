package kafka

import (
	"context"
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer    *kafka.Writer
	config    kafka.WriterConfig
	topic     string
	closeOnce sync.Once
	closeCh   chan struct{}
}

func NewProducer(config kafka.WriterConfig) (*Producer, error) {
	p := &Producer{
		config: config,
	}
	w := kafka.NewWriter(config)
	p.writer = w
	return p, nil
}

func (p *Producer) Send(ctx context.Context, msg kafka.Message) error {
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() {
	p.closeOnce.Do(func() {
		close(p.closeCh)
	})
}
