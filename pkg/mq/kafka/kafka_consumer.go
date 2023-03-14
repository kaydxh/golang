package kafka

import (
	"context"
	"sync"

	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	reader     *kafka.Reader
	config     kafka.ReaderConfig
	topic      string
	groupID    string
	msgChannel chan kafka.Message
	streamOnce sync.Once
	closeOnce  sync.Once
	closeCh    chan struct{}
}

func NewConsumer(config kafka.ReaderConfig, topic string, groupID string) (*Consumer, error) {
	c := &Consumer{
		config:     config,
		topic:      topic,
		groupID:    groupID,
		msgChannel: make(chan kafka.Message, 1024),
		closeCh:    make(chan struct{}),
	}
	r := kafka.NewReader(config)

	c.reader = r
	return c, nil
}

func (c *Consumer) ReadStream(ctx context.Context) <-chan kafka.Message {
	c.streamOnce.Do(func() {
		go func() {
			for {
				select {
				case <-c.closeCh:
					err := c.reader.Close()
					if err != nil {
						logrus.WithError(err).Errorf("failed to close consumer")
					}

					if c.msgChannel != nil {
						close(c.msgChannel)
					}

					return

				default:
					msg, err := c.reader.ReadMessage(ctx)
					if err != nil {
						logrus.WithError(err).Errorf("failed to read message")
						continue
					}
					c.msgChannel <- msg
				}
			}
		}()
	})

	return c.msgChannel
}

func (c *Consumer) Close() {
	c.closeOnce.Do(func() {
		close(c.closeCh)
	})
}
