package kafka_test

import (
	"testing"
	"time"

	kafka_ "github.com/kaydxh/golang/pkg/mq/kafka"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

func TestCreateTopic(t *testing.T) {
	ctx := context.Background()

	mq := kafka_.NewMQ(kafka_.MQConfig{
		Brokers: []string{"localhost:9092"},
	})
	_, err := mq.InstallMQ(ctx, time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("failed to new mq, err: %v", err)
	}
	topic := "topic-test-1"
	err = mq.CreateTopics(
		kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	if err != nil {
		t.Fatalf("failed to create topic, err: %v", err)
	}

}

func TestProducer(t *testing.T) {

	ctx := context.Background()
	mq := kafka_.NewMQ(kafka_.MQConfig{
		Brokers: []string{"localhost:9092"},
	})
	_, err := mq.InstallMQ(ctx, time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("failed to new mq, err: %v", err)
	}

	topic := "topic-test-1"
	err = mq.AsProducers(ctx, topic)
	if err != nil {
		t.Fatalf("failed to as producers, err: %v", err)
	}
	p, err := mq.GetProducer(topic)
	if err != nil {
		t.Fatalf("failed to get producer, err: %v", err)
	}

	err = p.Send(ctx,
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		t.Fatalf("failed to send messages, err: %v", err)
	}

}

func TestConsumer(t *testing.T) {

	ctx := context.Background()
	mq := kafka_.NewMQ(kafka_.MQConfig{
		Brokers: []string{"localhost:9092"},
	})
	_, err := mq.InstallMQ(ctx, time.Second, 10*time.Second)
	if err != nil {
		t.Fatalf("failed to new mq, err: %v", err)
	}

	topic := "topic-test-1"
	err = mq.AsConsumers(ctx, topic)
	if err != nil {
		t.Fatalf("failed to as producers, err: %v", err)
	}
	c, err := mq.GetConsumer(topic)
	if err != nil {
		t.Fatalf("failed to get producer, err: %v", err)
	}

	for msg := range c.ReadStream(ctx) {
		if msg.Error() != nil {
			t.Errorf("failed to read stream err: %v", err)
			continue
		}
		stas := c.Stats()
		t.Logf("read msg key[%v], value[%v], stas: %+v", string(msg.Key()), string(msg.Value()), stas)
	}

}

func TestNew(t *testing.T) {

	cfgFile := "./kafka.yaml"
	config := kafka_.NewConfig(kafka_.WithViper(viper_.GetViper(cfgFile, "mq.kafka")))

	mq, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("mq: %#v", mq)
}
