package kafka

import (
	"context"
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	*kafka.Writer
	config    kafka.WriterConfig
	topic     string
	closeOnce sync.Once
	closeCh   chan struct{}
}

func NewProducer(config kafka.WriterConfig) (*Producer, error) {
	p := &Producer{
		config:  config,
		closeCh: make(chan struct{}),
	}
	w := kafka.NewWriter(config)
	p.Writer = w
	return p, nil
}

func (p *Producer) Send(ctx context.Context, msgs ...kafka.Message) error {
	return p.Writer.WriteMessages(ctx, msgs...)
}

func (p *Producer) Close() {
	p.closeOnce.Do(func() {
		if p.closeCh != nil {
			close(p.closeCh)
		}
	})
}
