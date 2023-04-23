/*
 *Copyright (c) 2023, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package kafka

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
	time_ "github.com/kaydxh/golang/go/time"
	mq_ "github.com/kaydxh/golang/pkg/mq"
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
	// The default is to use a target batch size of 100 messages.
	batchSize int

	// Limit the maximum size of a request in bytes before being
	// sent to
	// a partition.
	//
	// The default is to use a kafka default value of
	// 1048576.
	batchBytes int

	// Time limit on how often incomplete message batches will be
	// flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	batchTimeout time.Duration
}

type ConsumerOptions struct {
	groupID   string
	partition int

	// MinBytes indicates to the broker the minimum batch size that the consumer
	// will accept. Setting a high minimum when consuming from a low-volume topic
	// may result in delayed delivery when the broker does not have enough data to
	// satisfy the defined minimum.
	//
	// Default: 1
	minBytes int

	// MaxBytes indicates to the broker the maximum batch size that the consumer
	// will accept. The broker will truncate a message to satisfy this maximum, so
	// choose a value that is high enough for your largest message size.
	//
	// Default: 1MB
	maxBytes int

	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	//
	// Default: 10s
	maxWait time.Duration
}

type MQ struct {
	*kafka.Conn
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

func newController(broker string) (*kafka.Conn, error) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, err
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return nil, err
	}

	return controllerConn, nil
}

func (q *MQ) InstallMQ(
	ctx context.Context,
	maxWaitInterval time.Duration,
	failAfter time.Duration,
) (*MQ, error) {
	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)

	var (
		errs []error
		conn *kafka.Conn
	)
	err := time_.BackOffUntilWithContext(ctx, func(ctx context.Context) (err_ error) {
		for _, broker := range q.Conf.Brokers {
			conn, err_ = newController(broker)
			if err_ != nil {
				errs = append(errs, err_)
				continue
			}
			return nil
		}
		return fmt.Errorf("failed to connect kafka: %v, err: %v", q.Conf.Brokers, errors_.NewAggregate(errs))
	}, exp, true, false)
	if err != nil {
		return nil, fmt.Errorf("get kafka connection fail after: %v", failAfter)
	}

	q.Conn = conn

	return q, nil
}

func (q *MQ) AsProducers(ctx context.Context, topics ...string) (producers []*Producer, err error) {
	for _, topic := range topics {

		fn := func() (*Producer, error) {

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
				return nil, err
			}

			q.producerLock.Lock()
			defer q.producerLock.Unlock()
			q.producers[topic] = producer
			return producer, nil
		}

		exp := time_.NewExponentialBackOff(
			time_.WithExponentialBackOffOptionMaxInterval(q.opts.reconnectBackOff),
			time_.WithExponentialBackOffOptionMaxElapsedTime(q.opts.reconnectBackOffMax),
		)
		err := time_.BackOffUntilWithContext(ctx, func(ctx context.Context) (err_ error) {
			producer, err_ := fn()
			if err_ != nil {
				return err_
			}
			producers = append(producers, producer)
			return nil
		}, exp, true, false)
		if err != nil {
			return nil, fmt.Errorf("create producer for %v fail after: %v", topic, q.opts.reconnectBackOffMax.Milliseconds())
		}

	}

	return producers, nil
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

func (q *MQ) Send(ctx context.Context, topic string, msgs ...kafka.Message) error {
	p, err := q.GetProducer(topic)
	if err != nil {
		return err
	}

	return p.Send(ctx, msgs...)
}

func (q *MQ) AsConsumers(ctx context.Context, topics ...string) (consumers []*Consumer, err error) {
	for _, topic := range topics {

		checkFn := func() bool {
			q.consumerLock.Lock()
			defer q.consumerLock.Unlock()
			_, ok := q.consumers[topic]
			if ok {
				return true
			}
			return false
		}

		exist := checkFn()
		if exist {
			continue
		}

		fn := func() (*Consumer, error) {

			dialer := &kafka.Dialer{
				Timeout:   q.opts.dialTimeout,
				DualStack: true,
			}
			consumer, err := NewConsumer(kafka.ReaderConfig{
				Brokers:  q.Conf.Brokers,
				GroupID:  q.opts.consumerOpts.groupID,
				Topic:    topic,
				Dialer:   dialer,
				MinBytes: q.opts.consumerOpts.minBytes,
				MaxBytes: q.opts.consumerOpts.maxBytes,
				MaxWait:  q.opts.consumerOpts.maxWait,
				//	CommitInterval: time.Second, // flushes commits to Kafka every second
			})
			if err != nil {
				return nil, err
			}

			q.consumerLock.Lock()
			defer q.consumerLock.Unlock()
			q.consumers[topic] = consumer
			return consumer, nil
		}

		exp := time_.NewExponentialBackOff(
			time_.WithExponentialBackOffOptionMaxInterval(q.opts.reconnectBackOff),
			time_.WithExponentialBackOffOptionMaxElapsedTime(q.opts.reconnectBackOffMax),
		)
		err := time_.BackOffUntilWithContext(ctx, func(ctx context.Context) (err_ error) {
			consumer, err_ := fn()
			if err_ != nil {
				return err_
			}
			consumers = append(consumers, consumer)
			return nil
		}, exp, true, false)
		if err != nil {
			return nil, fmt.Errorf("create consumer for %v fail after: %v", topic, q.opts.reconnectBackOffMax.Milliseconds())
		}

	}

	return consumers, nil
}

func (q *MQ) GetConsumer(topic string) (*Consumer, error) {
	q.consumerLock.Lock()
	defer q.consumerLock.Unlock()
	consumer, ok := q.consumers[topic]
	if ok {
		return consumer, nil
	}

	return nil, fmt.Errorf("not exist consumer %v", topic)
}

func (q *MQ) ReadStream(ctx context.Context, topic string) <-chan mq_.Message {
	c, err := q.GetConsumer(topic)
	if err != nil {
		return nil
	}

	return c.ReadStream(ctx)
}

func (q *MQ) Close() {
	for _, producer := range q.producers {
		producer.Close()
	}
	for _, consumer := range q.consumers {
		consumer.Close()
	}
	if q.Conn != nil {
		q.Conn.Close()
	}
}
