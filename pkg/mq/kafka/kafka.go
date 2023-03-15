package kafka

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
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
	*kafka.Conn
	producers map[string]*Producer // key topic
	consumers map[string]*Consumer // key topic

	producerLock sync.Mutex
	consumerLock sync.Mutex

	Conf MQConfig
	opts MQOptions
}

func NewMQ(conf MQConfig, opts ...MQOption) (*MQ, error) {
	c := &MQ{
		Conf:      conf,
		producers: make(map[string]*Producer),
		consumers: make(map[string]*Consumer),
	}
	c.ApplyOptions(opts...)

	var errs []error

	for _, broker := range conf.Brokers {
		conn, err := newController(broker)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		c.Conn = conn
		break
	}

	return c, errors_.NewAggregate(errs)

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

func (q *MQ) AsProducers(ctx context.Context, topics ...string) error {
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

func (q *MQ) Send(ctx context.Context, topic string, msgs ...kafka.Message) error {
	p, err := q.GetProducer(topic)
	if err != nil {
		return err
	}

	return p.Send(ctx, msgs...)
}

func (q *MQ) AsConsumers(ctx context.Context, topics ...string) error {
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

		fn := func() error {

			dialer := &kafka.Dialer{
				Timeout:   q.opts.dialTimeout,
				DualStack: true,
			}
			consumer, err := NewConsumer(kafka.ReaderConfig{
				Brokers: q.Conf.Brokers,
				//GroupID: "",
				Topic:  topic,
				Dialer: dialer,
				//	MinBytes:       10e3,        // 10KB
				//	MaxBytes:       10e6,        // 10MB
				//	CommitInterval: time.Second, // flushes commits to Kafka every second
			})
			if err != nil {
				return err
			}

			q.consumerLock.Lock()
			defer q.consumerLock.Unlock()
			q.consumers[topic] = consumer
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
			return fmt.Errorf("create consumer for %v fail after: %v", topic, q.opts.reconnectBackOffMax.Milliseconds())
		}

	}

	return nil
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

func (q *MQ) ReadStream(ctx context.Context, topic string) <-chan kafka.Message {
	c, err := q.GetConsumer(topic)
	if err != nil {
		return nil
	}

	return c.ReadStream(ctx)
}

func (q *MQ) Close() {
	if q.Conn != nil {
		q.Conn.Close()
	}

	for _, producer := range q.producers {
		producer.Close()
	}
	for _, consumer := range q.consumers {
		consumer.Close()
	}
}
