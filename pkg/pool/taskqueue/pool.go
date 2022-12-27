package taskqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	"github.com/sirupsen/logrus"
)

type PoolOptions struct {
	workerBurst   uint32
	fetcherBurst  uint32
	fetchTimeout  time.Duration
	resultExpired time.Duration

	resultCallbackFunc queue_.ResultCallbackFunc
}

type Pool struct {
	msgChan chan *queue_.Message
	opts    PoolOptions
	taskq   queue_.Queue

	wg sync.WaitGroup
}

func defaultPoolOptions() PoolOptions {
	return PoolOptions{
		workerBurst:   1,
		fetcherBurst:  1,
		fetchTimeout:  10 * time.Second,
		resultExpired: 30 * time.Minute,
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

func (p *Pool) Publish(ctx context.Context, msg *queue_.Message) (string, error) {
	taskId, err := p.taskq.Add(ctx, msg)
	if err != nil {
		logrus.WithError(err).Errorf("failed to publish msg: %v", msg)
		return "", err
	}
	logrus.Infof("succeed in publishing msg: %v", msg)
	return taskId, nil
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

func (p *Pool) FetchResult(ctx context.Context, key string) (*queue_.MessageResult, error) {
	result, err := p.taskq.FetchResult(ctx, key)
	if err != nil {
		logrus.WithError(err).Errorf("failed to fetch result of msg %v", key)
		return nil, err
	}

	return result, nil
}

func (p *Pool) Process(ctx context.Context, msg *queue_.Message) error {

	tasker := Get(msg.Scheme)
	if tasker == nil {
		return fmt.Errorf("not regiter task scheme %v", msg.Scheme)
	}

	var errs []error
	result := &queue_.MessageResult{
		Id:      msg.Id,
		InnerId: msg.InnerId,
		Name:    msg.Name,
		Scheme:  msg.Scheme,
	}

	clean := func() {
		err := p.taskq.Delete(ctx, msg)
		if err != nil {
			// only log error
			logrus.WithError(err).Errorf("failed to delete msg: %v", msg)
			return
		}
		logrus.Infof("delete msg: %v", msg)
	}
	defer clean()

	result, err := tasker.TaskHandler(ctx, msg)
	if err != nil {
		errs = append(errs, err)
		logrus.WithError(err).Errorf("failed to handle task %v, err: %v", msg, err)
	}
	result.Err = err

	//callback result
	if p.opts.resultCallbackFunc != nil {
		p.opts.resultCallbackFunc(ctx, result)
	}

	//write to queue
	_, err = p.taskq.AddResult(ctx, result, p.opts.resultExpired)
	if err != nil {
		errs = append(errs, err)
		//only log error
		logrus.WithError(err).Errorf("failed to add msg result: %v", result)
		return errors_.NewAggregate(errs)
	}

	return nil
}
