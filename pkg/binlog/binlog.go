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
package binlog

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaydxh/golang/go/errors"
	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
	mq_ "github.com/kaydxh/golang/pkg/mq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BinlogOptions struct {
	rootPath       string
	prefixName     string
	suffixName     string
	flushBatchSize int
	flushInterval  time.Duration
	rotateInterval time.Duration
	rotateSize     int64
}

type Channel struct {
	Name string
}

type BinlogService struct {
	consumers []mq_.Consumer

	// 所分配的channel
	//channels    []Channel
	rotateFiler *rotate_.RotateFiler

	opts BinlogOptions

	inShutdown atomic.Bool
	mu         sync.Mutex
	cancel     func()
}

func defaultBinlogServiceOptions() BinlogOptions {
	opts := BinlogOptions{
		prefixName:     "segment",
		suffixName:     "log",
		flushBatchSize: 1024,
		flushInterval:  time.Second, // 1s
		rotateInterval: time.Hour,
		rotateSize:     512 * 1024 * 1024, //512M
	}
	path, err := os.Getwd()
	if err != nil {
		path = "/"
	}
	opts.rootPath = path
	return opts
}

func NewBinlogService(consumers []mq_.Consumer, opts ...BinlogServiceOption) (*BinlogService, error) {
	if len(consumers) == 0 {
		return nil, fmt.Errorf("consumers is empty")
	}

	bs := &BinlogService{
		consumers: consumers,
		opts:      defaultBinlogServiceOptions(),
	}
	bs.ApplyOptions(opts...)

	rotateFiler, _ := rotate_.NewRotateFiler(
		bs.opts.rootPath,
		rotate_.WithRotateSize(bs.opts.rotateSize),
		rotate_.WithRotateInterval(bs.opts.rotateInterval),
		rotate_.WithSuffixName(bs.opts.suffixName),
		rotate_.WithPrefixName(bs.opts.prefixName),
	)
	bs.rotateFiler = rotateFiler

	return bs, nil
}

func (srv *BinlogService) logger() logrus.FieldLogger {
	return logrus.WithField("module", "BinlogService")
}

func (srv *BinlogService) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("BinlogService Run")
	if srv.inShutdown.Load() {
		logger.Infoln("BinlogService Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil

}

func (srv *BinlogService) flush(ctx context.Context, consumer mq_.Consumer) error {
	logger := srv.logger()
	timer := time.NewTicker(srv.opts.flushInterval)
	defer timer.Stop()

	var flushBatchData [][]byte
	for msg := range consumer.ReadStream(ctx) {
		if msg.Error() != nil {
			logger.WithError(msg.Error()).Errorf("faild to read stream %v", consumer.Channel())
			continue
		}

		flushFunc := func(ctx context.Context, data []byte) (err error) {
			flushBatchData = append(flushBatchData, data)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				if len(flushBatchData) > 0 {
					_, _, err = srv.rotateFiler.WriteBytesLine(flushBatchData)
					flushBatchData = nil
				}
			default:
				if len(flushBatchData) >= srv.opts.flushBatchSize {
					_, _, err = srv.rotateFiler.WriteBytesLine(flushBatchData)
					flushBatchData = nil
					return err
				}
				if err != nil {
					return err
				}
			}
			return nil
		}
		err := flushFunc(ctx, msg.Value())
		if err != nil {
			logger.WithError(msg.Error()).Errorf("faild to flush message for channel[%v]", consumer.Channel())
			return err
		}
	}

	return nil
}

func (srv *BinlogService) Serve(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ServiceBinlog Serve")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("ServiceBinlog Serve canceled")
		return err
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	g, gCtx := errgroup.WithContext(ctx)
	for _, consumer := range srv.consumers {
		consumer := consumer
		g.Go(func() error {
			return srv.flush(gCtx, consumer)
		})
	}
	err := g.Wait()
	if err != nil {
		logger.WithError(err).Errorf("wait flush worker")
		return err
	}
	logger.Info("stopped binlog service")
	return nil
}

func (srv *BinlogService) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, consumer := range srv.consumers {
		if consumer != nil {
			consumer.Close()
		}
	}
	if srv.cancel != nil {
		srv.cancel()
	}
}
