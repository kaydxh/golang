package redisq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	os_ "github.com/kaydxh/golang/go/os"
	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	opts queue_.QueueOptions

	stream       string
	streamGroup  string
	consumerName string
	db           *redis.Client
}

func NewQueue(db *redis.Client, opts queue_.QueueOptions) *Queue {
	q := &Queue{
		db:           db,
		stream:       fmt.Sprintf("taskq-%s-stream", opts.Name),
		streamGroup:  "taskq",
		consumerName: os_.GetProcId(),
	}

	return q
}

// https://redis.io/commands/xadd/
// * 为自动生成Id
// XADD taskq-redis-stream  *  taskq "123"
func (q *Queue) Add(ctx context.Context, msg *queue_.Message) error {
	if msg.Name == "" {
		return fmt.Errorf("messgae name is empty")
	}
	if msg.Id == "" {
		msg.Id = uuid.NewString()
	}

	//todo check message name if exists
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// xadd msg to redis
	err = q.db.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: map[string]interface{}{
			"message": string(data),
		},
	}).Err()
	if err != nil {
		logrus.WithError(err).Errorf("failed to xadd msg [name: %v, id: %v]", msg.Name, msg.Id)
		return err
	}

	return nil
}

// get all
// XRANGE taskq-redis-stream  - +

func (q *Queue) fetchN(ctx context.Context, n int64, waitTimeout time.Duration) ([]*queue_.Message, error) {
	streams, err := q.db.XReadGroup(ctx, &redis.XReadGroupArgs{
		Streams:  []string{q.stream, ">"},
		Group:    q.streamGroup,
		Consumer: q.consumerName,
		Count:    n,
		Block:    waitTimeout,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	stream := &streams[0]
	msgs := make([]*queue_.Message, len(stream.Messages))
	for i := range stream.Messages {
		xmsg := &stream.Messages[i]
		msgs[i] = &queue_.Message{}
		msg := msgs[i]
		//	msg.Ctx = ctx
		data, ok := xmsg.Values["message"].(string)
		if !ok {
			logrus.Errorf("message type is not string")
			continue
		}
		err = json.Unmarshal([]byte(data), msg)
		if err != nil {
			//msg.Err = err
			logrus.WithError(err).Errorf("failed to unmarshal msg %v", string(data))
			return nil, err
		}
		msg.InnerId = xmsg.ID
	}

	return msgs, nil
}

func (q *Queue) FetchN(ctx context.Context, n int64, waitTimeout time.Duration) ([]*queue_.Message, error) {
	msgs, err := q.fetchN(ctx, n, waitTimeout)
	if err != nil {
		if strings.HasPrefix(err.Error(), "NOGROUP") {
			q.createMkStreamGroup(ctx)
			return q.fetchN(ctx, n, waitTimeout)
		}
	}
	return msgs, err
}

func (q *Queue) FetchOne(ctx context.Context, waitTimeout time.Duration) (*queue_.Message, error) {
	msgs, err := q.FetchN(ctx, 1, waitTimeout)
	if err != nil {
		return nil, err
	}

	if len(msgs) == 0 {
		//return nil, fmt.Errorf("fetch %v number messages", len(msgs))
		return nil, nil
	}
	return msgs[0], nil
}

func (q *Queue) createMkStreamGroup(ctx context.Context) {
	logrus.Infof("createMkStreamGroup stream: %v, streamGroup: %v", q.stream, q.streamGroup)
	_ = q.db.XGroupCreateMkStream(ctx, q.stream, q.streamGroup, "0").Err()
}

// XLEN taskq-redis-stream
func (q *Queue) Len() (int64, error) {
	return q.db.XLen(context.Background(), q.stream).Result()
}

// XACK taskq-redis-stream taskq 1671960330664-0
// XDEL taskq-redis-stream 1671960330664-0
func (q *Queue) Delete(ctx context.Context, msg *queue_.Message) error {
	err := q.db.XAck(ctx, q.stream, q.streamGroup, msg.InnerId).Err()
	if err != nil {
		logrus.WithError(err).Errorf("failed to XAck msg %v, stream[%v], stremaGroup[%v], id[%v], innerId[%v]", msg, q.stream, q.streamGroup, msg.Id, msg.InnerId)
		return err
	}
	return q.db.XDel(ctx, q.stream, msg.InnerId).Err()
}
