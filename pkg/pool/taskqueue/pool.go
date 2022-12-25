package taskqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	"github.com/sirupsen/logrus"
)

type PoolOptions struct {
	workerBurst  uint32
	fetcherBurst uint32
	fetchTimeout time.Duration
}

type Pool struct {
	msgChan chan *queue_.Message
	opts    PoolOptions
	taskq   queue_.Queue

	wg sync.WaitGroup
}

func defaultPoolOptions() PoolOptions {
	return PoolOptions{
		workerBurst:  1,
		fetcherBurst: 1,
		fetchTimeout: 10 * time.Second,
	}
}

func NewPool(taskq queue_.Queue, opts ...PoolOption) *Pool {
	p := &Pool{
		msgChan: make(chan *queue_.Message),
		taskq:   taskq,
		opts:    defaultPoolOptions(),
	}
	p.ApplyOptions(opts...)

	return p
}

func (p *Pool) Publish(ctx context.Context, msg *queue_.Message) error {
	err := p.taskq.Add(ctx, msg)
	if err != nil {
		logrus.WithError(err).Errorf("failed to publish msg: %v", msg)
		return err
	}
	logrus.Infof("succeed in publishing msg: %v", msg)
	return nil
}

func (p *Pool) Consume(ctx context.Context) error {

	var i uint32
	for i = 0; i < p.opts.workerBurst; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()

			for {
				select {
				case msg, ok := <-p.msgChan:
					if !ok {
						logrus.Errorf("message channel is closed.")
						return
					}
					err := p.Process(ctx, msg)
					if err != nil {
						logrus.Errorf("faild to process msg, err: %v", err)
					}
				case <-ctx.Done():
					//  err: context canceled
					return

				}
			}

		}()
	}

	for i = 0; i < p.opts.fetcherBurst; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()

			for {
				msg, err := p.taskq.FetchOne(ctx, p.opts.fetchTimeout)
				if err != nil {
					logrus.Errorf("faild to fetch msg, err: %v", err)
					continue
				}
				if msg == nil {
					logrus.Infof("no msg to fetch")
					continue
				}

				select {
				case p.msgChan <- msg:
				case <-ctx.Done():
					//  err: context canceled
					return

				}
			}

		}()

	}

	return nil
}

func (p *Pool) Process(ctx context.Context, msg *queue_.Message) error {

	tasker := Get(msg.Scheme)
	if tasker == nil {
		return fmt.Errorf("not regiter task scheme %v", msg.Scheme)
	}

	err := tasker.TaskHandler(msg)
	if err != nil {
		return fmt.Errorf("failed to handle task %v, err: %v", msg, err)
	}

	err = p.taskq.Delete(ctx, msg)
	if err != nil {
		// only log error
		logrus.WithError(err).Errorf("failed to delete msg: %v", msg)
	}

	return nil
}
