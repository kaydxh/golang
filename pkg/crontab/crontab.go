package crontab

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kaydxh/golang/go/errors"
	errors_ "github.com/kaydxh/golang/go/errors"
	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

// Crontab ...
type CrontabSerivce struct {
	checkInterval time.Duration
	inShutdown    atomic.Bool // true when when server is in shutdown

	mu     sync.Mutex
	cancel func()
	fs     []func(context.Context, *logrus.Entry) error
}

// NewCrontab ...
func NewCrontabSerivce(
	checkInterval time.Duration,
) *CrontabSerivce {
	s := &CrontabSerivce{
		checkInterval: checkInterval,
	}
	return s
}

// Register ...
func (c *CrontabSerivce) Register(f func(context.Context, *logrus.Entry) error) {
	c.fs = append(c.fs, f)
}

// Run will initialize the backend. It must not block, but may run go routines in the background.
func (c *CrontabSerivce) Run(ctx context.Context) error {
	logger := c.getLogger()
	logger.Infoln("Crontab Run")
	if c.inShutdown.Load() {
		logger.Infoln("Crontab Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(c.Serve(ctx))
	}()
	return nil
}

// Serve ...
func (c *CrontabSerivce) Serve(ctx context.Context) error {
	logger := c.getLogger()
	logger.Infoln("Crontab Serve")

	if c.inShutdown.Load() {
		logger.Infoln("Crontab Shutdown")
		return fmt.Errorf("server closed")
	}

	defer c.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c.mu.Lock()
	c.cancel = cancel
	c.mu.Unlock()

	time_.UntilWithContxt(ctx, func(ctx context.Context) {
		err := c.check(ctx)
		if err != nil {
			return
		}
	}, c.checkInterval)
	if err := ctx.Err(); err != nil {
		logger.WithError(err).Errorf("stopped checking")
		return err
	}
	logger.Info("stopped checking")
	return nil
}

// Shutdown ...
func (c *CrontabSerivce) Shutdown() {
	c.inShutdown.Store(true)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *CrontabSerivce) getLogger() *logrus.Entry {
	return logrus.WithField("module", "Crontab")
}

func (c *CrontabSerivce) check(ctx context.Context) error {
	var (
		wg   sync.WaitGroup
		errs []error
	)
	logger := c.getLogger()
	c.mu.Lock()
	fs := c.fs
	c.mu.Unlock()
	go func() {
		for _, f := range fs {
			err := f(ctx, logger)
			if err != nil {
				c.mu.Lock()
				errs = append(errs, err)
				c.mu.Unlock()
			}
		}
	}()
	wg.Wait()
	return errors_.NewAggregate(errs)
}

/*
type CronProcessor struct {
	cronRunner *cron.Cron
}

type Job interface {
	Run()
}

func NewCronProcessor(specTime string, job Job) (*CronProcessor, error) {
	cp := &CronProcessor{
		cronRunner: cron.New(),
	}

	if cp.cronRunner == nil {
		return nil, fmt.Errorf("failed to init cronRunner")
	}

	if specTime == "" {
		return nil, fmt.Errorf("specTime is empty")
	}

	err := cp.cronRunner.AddJob(specTime, job)
	if err != nil {
		return nil, fmt.Errorf("failed to odd job in specTime: [%v]", specTime)
	}

	cp.cronRunner.Start()

	return cp, nil
}

func (c *CronProcessor) Stop() {
	c.cronRunner.Stop()
}
*/