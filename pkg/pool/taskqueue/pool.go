package taskqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
	time_ "github.com/kaydxh/golang/go/time"
	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	"github.com/sirupsen/logrus"
)

type PoolOptions struct {
	workerBurst    uint32
	fetcherBurst   uint32
	fetchTimeout   time.Duration
	resultExpired  time.Duration
	processTimeout time.Duration

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

	result := &queue_.MessageResult{
		Id:      msg.Id,
		InnerId: msg.InnerId,
		Name:    msg.Name,
		Scheme:  msg.Scheme,
		Status:  queue_.MessageStatus_Doing,
	}

	//write to queue
	_, err = p.taskq.AddResult(ctx, result, p.opts.resultExpired)
	if err != nil {
		//only log error, this error not stop to handle msg
		logrus.WithError(err).Errorf("failed to add msg result: %v", result)
	}
	return taskId, nil
}

func (p *Pool) Consume(ctx context.Context) (err error) {

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

				const backoff = time.Second
				msg, err := p.taskq.FetchOne(ctx, p.opts.fetchTimeout)
				if err != nil {
					logrus.Errorf("faild to fetch msg, err: %v", err)
					//todo backoff
					time.Sleep(backoff)
					continue
				}
				if msg == nil {
					logrus.Debugf("no msg to fetch")
					time.Sleep(backoff)
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

	result := &queue_.MessageResult{
		Id:      msg.Id,
		InnerId: msg.InnerId,
		Name:    msg.Name,
		Scheme:  msg.Scheme,
		Status:  queue_.MessageStatus_Doing,
	}
	err := time_.CallWithTimeout(ctx, p.opts.processTimeout, func(ctx context.Context) error {
		result_, err_ := tasker.TaskHandler(ctx, msg)
		if err_ != nil {
			errs = append(errs, err_)
			logrus.WithError(err_).Errorf("failed to handle task %v, err: %v", msg, err_)
		} else {
			if result_ != nil {
				result = result_
			}
		}
		return err_
	})
	result.Err = err

	if err == time_.ErrTimeout || err == context.Canceled {
		result.Status = queue_.MessageStatus_Fail
	} else {
		// success means the message is processed, but not means hander
		// function return nil, just not run timeout or canceled
		result.Status = queue_.MessageStatus_Success
	}

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
