package kafka

import (
	"fmt"
	"sync"
	"time"

	time_ "github.com/kaydxh/golang/go/time"
	kafka "github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type MQConfig struct {
	Brokers []string
}

type MQOptions struct {
	SaslUsername string
	SaslPassword string

	dialTimeout         time.Duration
	reconnectBackOff    time.Duration
	reconnectBackOffMax time.Duration

	producerOpts ProducerOptions
	consumerOpts ConsumerOptions
}

type ProducerOptions struct {
}

type ConsumerOptions struct {
}

type MQ struct {
	producers map[string]*Producer // key topic
	consumers map[string]*Consumer // key topic

	producerLock sync.Mutex
	consumerLock sync.Mutex

	Conf MQConfig
	opts MQOptions
}

func NewMQ(conf MQConfig, opts ...MQOption) *MQ {
	c := &MQ{
		Conf:      conf,
		producers: make(map[string]*Producer),
		consumers: make(map[string]*Consumer),
	}

	c.ApplyOptions(opts...)
	return c
}

func (q *MQ) AsProducers(ctx context.Context, topics []string) error {
	for _, topic := range topics {

		fn := func() error {

			dialer := &kafka.Dialer{
				Timeout:   q.opts.dialTimeout,
				DualStack: true,
			}
			producer, err := NewProducer(kafka.WriterConfig{
				Brokers:  q.Conf.Brokers,
				Topic:    topic,
				Balancer: &kafka.Hash{},
				Dialer:   dialer,
			})
			if err != nil {
				return err
			}

			q.producerLock.Lock()
			defer q.producerLock.Unlock()
			q.producers[topic] = producer
			return nil
		}

		exp := time_.NewExponentialBackOff(
			time_.WithExponentialBackOffOptionMaxInterval(q.opts.reconnectBackOff),
			time_.WithExponentialBackOffOptionMaxElapsedTime(q.opts.reconnectBackOffMax),
		)
		err := time_.BackOffUntilWithContext(ctx, func(ctx context.Context) (err_ error) {
			err_ = fn()
			if err_ != nil {
				return err_
			}
			return nil
		}, exp, true, false)
		if err != nil {
			return fmt.Errorf("create producer for %v fail after: %v", topic, q.opts.reconnectBackOffMax.Milliseconds())
		}

	}

	return nil
}

func (q *MQ) GetProducer(topic string) (*Producer, error) {
	q.producerLock.Lock()
	defer q.producerLock.Unlock()
	producer, ok := q.producers[topic]
	if ok {
		return producer, nil
	}

	return nil, fmt.Errorf("not exist producer %v", topic)
}
