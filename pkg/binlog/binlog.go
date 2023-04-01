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
	maxFlushBatchSize int
	flushInterval     time.Duration
	maxRotateInterval time.Duration
	maxRotateSize     int64
}

type Channel struct {
	Name string
}

type BinlogService struct {
	consumer mq_.Consumer

	// 所分配的channel
	channels    []Channel
	rotateFiler *rotate_.RotateFiler

	opts BinlogOptions

	inShutdown atomic.Bool
	mu         sync.Mutex
	cancel     func()
}

func defaultBinlogServiceOptions() BinlogOptions {
	return BinlogOptions{
		maxFlushBatchSize: 1024,
		flushInterval:     time.Second, // 1s
		maxRotateInterval: time.Hour,
		maxRotateSize:     512 * 1024 * 1024, //512M
	}
}

func NewBinlogService(channels []Channel, opts ...BinlogServiceOption) *BinlogService {
	bs := &BinlogService{
		channels: channels,
		opts:     defaultBinlogServiceOptions(),
	}
	bs.ApplyOptions(opts...)
	return bs
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

func (srv *BinlogService) flush(ctx context.Context, channel Channel) error {
	logger := srv.logger()
	timer := time.NewTicker(srv.opts.flushInterval)
	defer timer.Stop()

	var flushBatchData [][]byte
	for msg := range srv.consumer.ReadStream(ctx, channel.Name) {
		if msg.Error() != nil {
			logger.WithError(msg.Error()).Errorf("faild to read stream %v", channel)
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
				if len(flushBatchData) >= srv.opts.maxFlushBatchSize {
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
			logger.WithError(msg.Error()).Errorf("faild to flush message for channel[%v]", channel)
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
	for _, channel := range srv.channels {
		channel := channel
		g.Go(func() error {
			return srv.flush(gCtx, channel)
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
	if srv.cancel != nil {
		srv.cancel()
	}
}
