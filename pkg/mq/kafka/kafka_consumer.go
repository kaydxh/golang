package kafka

import (
	"context"
	"sync"

	mq_ "github.com/kaydxh/golang/pkg/mq"
	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	*kafka.Reader
	config     kafka.ReaderConfig
	msgChannel chan mq_.Message
	streamOnce sync.Once
	closeOnce  sync.Once
	closeCh    chan struct{}
}

func NewConsumer(config kafka.ReaderConfig) (*Consumer, error) {
	c := &Consumer{
		config:     config,
		msgChannel: make(chan mq_.Message, 1024),
		closeCh:    make(chan struct{}),
	}
	r := kafka.NewReader(config)

	c.Reader = r
	return c, nil
}

func (c *Consumer) Topic() string {
	return c.config.Topic
}

func (c *Consumer) ReadStream(ctx context.Context) <-chan mq_.Message {
	c.streamOnce.Do(func() {
		go func() {
			for {
				select {
				case <-c.closeCh:
					err := c.Reader.Close()
					if err != nil {
						logrus.WithError(err).Errorf("failed to close consumer")
					}

					if c.msgChannel != nil {
						close(c.msgChannel)
					}

					return

				default:
					msg, err := c.Reader.ReadMessage(ctx)
					c.msgChannel <- KafkaMessage{
						Err: err,
						Msg: &msg,
					}
				}
			}
		}()
	})

	return c.msgChannel
}

func (c *Consumer) Close() {
	c.closeOnce.Do(func() {
		if c.closeCh != nil {
			close(c.closeCh)
		}
	})
}
